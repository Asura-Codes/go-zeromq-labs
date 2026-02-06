package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"gemini-zeromq-labs/lab12/internal/config"
	"gemini-zeromq-labs/lab12/internal/mdp"
	"github.com/go-zeromq/zmq4"
)

func main() {
	name := flag.String("name", "Client", "Client name for logging")
	flag.Parse()

	log.Printf("[%s] Starting Client Requester...", *name)

	ctx := context.Background()
	client := zmq4.NewReq(ctx)
	defer client.Close()

	if err := client.Dial(config.BrokerFrontendAddr); err != nil {
		log.Fatalf("[%s] Failed to dial broker: %v", *name, err)
	}

	for i := 1; i <= 5; i++ {
		reqBody := fmt.Sprintf("Hello Majordomo %d from %s", i, *name)
		log.Printf("[%s] Sending: %s", *name, reqBody)

		// REQ socket automatically adds an empty frame before the message?
		// No, REQ/REP pair handles the empty delimiter frame internally or at the boundary.
		// However, when talking to a ROUTER (Broker), the Broker sees [Identity][Empty][Body].
		// The Client using REQ sends [Body].
		// But MDPC01 requires: [MDPC01][REQUEST][Service][Body].
		// So we must put that in the body.
		
		// Wait, MDPC01 spec says: [Empty][MDPC01][REQUEST]...
		// Does REQ add the Empty frame?
		// If Client is REQ and Broker is ROUTER:
		// Client Send("ABC") -> Broker Recv([Identity][Empty][ABC])
		// So the "Empty" frame required by MDP is likely the one REQ provides?
		// "MDPC01" spec assumes DEALER/ROUTER?
		// Spec: "The client connects... using a REQ socket."
		// "The client sends a request... [MDPC01][REQUEST][Service][Body]"
		// The "Empty" frame is implied by the REQ/ROUTER connection? 
		// Actually, if we send [MDPC01]... via REQ, the Broker (ROUTER) receives [ID][Empty][MDPC01]...
		// So yes, we just send the content frames.

		frames := [][]byte{
			[]byte(mdp.ClientHeader),
			[]byte(mdp.ClientRequest),
			[]byte(config.ServiceNameEcho),
			[]byte(reqBody),
		}
		
		err := client.Send(zmq4.Msg{Frames: frames})
		if err != nil {
			log.Fatalf("Send error: %v", err)
		}

		// Receive Reply
		reply, err := client.Recv()
		if err != nil {
			log.Fatalf("Recv error: %v", err)
		}

		log.Printf("Debug: Received %d frames", len(reply.Frames))
		for idx, f := range reply.Frames {
			log.Printf("  Frame %d: %q", idx, f)
		}

		// Reply: [MDPC01][REPLY][Service][Body] (Broker logic I wrote)
		// Check header
		if len(reply.Frames) >= 3 {
			log.Printf("[%s] Received: %s", *name, string(reply.Frames[len(reply.Frames)-1]))
		} else {
			log.Printf("[%s] Received malformed: %v", *name, reply.Frames)
		}
		
		time.Sleep(1 * time.Second)
	}
}
