package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-zeromq/zmq4"
	"github.com/google/uuid"
)

func main() {
	addr := flag.String("addr", "tcp://127.0.0.1:5557", "Address to listen on")
	flag.Parse()

	log.Printf("Starting Storage Service on %s...", *addr)

	// Create storage directories
	storageDir := "titanic_data"
	queueDir := filepath.Join(storageDir, "queue")
	replyDir := filepath.Join(storageDir, "replies")

	os.MkdirAll(queueDir, 0755)
	os.MkdirAll(replyDir, 0755)

	ctx := context.Background()
	rep := zmq4.NewRep(ctx)
	defer rep.Close()

	if err := rep.Listen(*addr); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	for {
		msg, err := rep.Recv()
		if err != nil {
			log.Printf("Recv error: %v", err)
			continue
		}

		frames := msg.Frames
		if len(frames) < 1 {
			continue
		}

		command := string(frames[0])
		var reply [][]byte

		switch command {
		case "SAVE":
			// SAVE [Service] [Data...]
			if len(frames) < 3 {
				reply = [][]byte{[]byte("ERROR"), []byte("Missing data")}
			} else {
				id := uuid.New().String()
				service := string(frames[1])
				data := frames[2:]

				// Save to queue: filename = ID.service
				filename := fmt.Sprintf("%s.%s", id, service)
				path := filepath.Join(queueDir, filename)

				// Flatten data for simple storage (or use multi-part format)
				// For simplicity, we just save the first frame of data
				err := ioutil.WriteFile(path, data[0], 0644)
				if err != nil {
					reply = [][]byte{[]byte("ERROR"), []byte(err.Error())}
				} else {
					reply = [][]byte{[]byte("OK"), []byte(id)}
					log.Printf("Stored request %s for %s", id, service)
				}
			}

		case "LIST":
			// LIST pending requests
			files, _ := ioutil.ReadDir(queueDir)
			reply = [][]byte{[]byte("OK")}
			for _, f := range files {
				reply = append(reply, []byte(f.Name()))
			}

		case "GET":
			// GET [filename]
			if len(frames) < 2 {
				reply = [][]byte{[]byte("ERROR")}
			} else {
				filename := string(frames[1])
				path := filepath.Join(queueDir, filename)
				content, err := ioutil.ReadFile(path)
				if err != nil {
					reply = [][]byte{[]byte("ERROR")}
				} else {
					reply = [][]byte{[]byte("OK"), content}
				}
			}

		case "COMPLETE":
			// COMPLETE [ID.service] [ReplyData]
			if len(frames) < 3 {
				reply = [][]byte{[]byte("ERROR")}
			} else {
				filename := string(frames[1])
				id := filename[:36] // Extract UUID
				data := frames[2]

				// Move from queue to replies
				oldPath := filepath.Join(queueDir, filename)
				newPath := filepath.Join(replyDir, id)

				err := ioutil.WriteFile(newPath, data, 0644)
				if err == nil {
					os.Remove(oldPath)
					reply = [][]byte{[]byte("OK")}
					log.Printf("Completed request %s", id)
				} else {
					reply = [][]byte{[]byte("ERROR")}
				}
			}

		case "FETCH":
			// FETCH [ID]
			if len(frames) < 2 {
				reply = [][]byte{[]byte("ERROR")}
			} else {
				id := string(frames[1])
				path := filepath.Join(replyDir, id)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					// Check if still in queue
					inQueue := false
					files, _ := ioutil.ReadDir(queueDir)
					for _, f := range files {
						if len(f.Name()) >= 36 && f.Name()[:36] == id {
							inQueue = true
							break
						}
					}
					if inQueue {
						reply = [][]byte{[]byte("PENDING")}
					} else {
						reply = [][]byte{[]byte("UNKNOWN")}
					}
				} else {
					content, _ := ioutil.ReadFile(path)
					reply = [][]byte{[]byte("OK"), content}
				}
			}

		case "DELETE":
			// DELETE [ID]
			if len(frames) < 2 {
				reply = [][]byte{[]byte("ERROR")}
			} else {
				id := string(frames[1])
				path := filepath.Join(replyDir, id)
				err := os.Remove(path)
				if err == nil {
					reply = [][]byte{[]byte("OK")}
					log.Printf("Deleted request %s after client acknowledgement", id)
				} else {
					reply = [][]byte{[]byte("ERROR"), []byte(err.Error())}
				}
			}
		}

		rep.Send(zmq4.Msg{Frames: reply})
	}
}
