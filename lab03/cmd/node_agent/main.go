package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gemini-zeromq-labs/lab03/internal/config"
	"gemini-zeromq-labs/lab03/internal/protocol"

	"github.com/go-zeromq/zmq4"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	cfg := config.LoadConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("Shutdown signal received")
		cancel()
	}()

	// Create REP socket
	socket := zmq4.NewRep(ctx)
	defer socket.Close()

	bindAddr := cfg.BindAddr()
	logger.Info("Node Agent listening", "endpoint", bindAddr)
	if err := socket.Listen(bindAddr); err != nil {
		logger.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	for {
		// Receive Request
		msg, err := socket.Recv()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			logger.Error("Error receiving message", "error", err)
			continue
		}

		logger.Info("Received request", "bytes", len(msg.Bytes()))

		// Parse Request
		req, err := protocol.FromBytes(msg.Bytes())
		if err != nil {
			sendError(socket, "Invalid JSON format")
			continue
		}

		// Process Command
		response := processCommand(req, logger)

		// Send Response
		respBytes, _ := response.ToBytes()
		if err := socket.Send(zmq4.NewMsg(respBytes)); err != nil {
			logger.Error("Error sending response", "error", err)
		}
	}
}

func processCommand(req *protocol.Request, logger *slog.Logger) protocol.Response {
	var data interface{}
	var err error

	logger.Info("Processing command", "command", req.Command)

	switch req.Command {
	case protocol.CMD_CPU:
		data, err = getCPUInfo()
	case protocol.CMD_MEM:
		data, err = getMemInfo()
	case protocol.CMD_HOST:
		data, err = getHostInfo()
	default:
		return protocol.Response{Status: "ERROR", Error: "Unknown command"}
	}

	if err != nil {
		logger.Error("Error collecting metrics", "command", req.Command, "error", err)
		return protocol.Response{Status: "ERROR", Error: err.Error()}
	}

	return protocol.Response{Status: "OK", Data: data}
}

func sendError(socket zmq4.Socket, msg string) {
	resp := protocol.Response{Status: "ERROR", Error: msg}
	bytes, _ := resp.ToBytes()
	socket.Send(zmq4.NewMsg(bytes))
}

// System Metric Helpers

func getCPUInfo() (protocol.CPUData, error) {
	info, err := cpu.Info()
	if err != nil {
		return protocol.CPUData{}, err
	}
	percent, err := cpu.Percent(time.Second, false) // 1 second sample
	if err != nil {
		return protocol.CPUData{}, err
	}

	model := "Unknown"
	if len(info) > 0 {
		model = info[0].ModelName
	}
	
	cores, _ := cpu.Counts(true)

	return protocol.CPUData{
		Model:        model,
		Cores:        cores,
		UsagePercent: percent,
	}, nil
}

func getMemInfo() (protocol.MemData, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return protocol.MemData{}, err
	}
	return protocol.MemData{
		Total:       v.Total,
		Available:   v.Available,
		UsedPercent: v.UsedPercent,
	}, nil
}

func getHostInfo() (protocol.HostData, error) {
	h, err := host.Info()
	if err != nil {
		return protocol.HostData{}, err
	}
	return protocol.HostData{
		Hostname: h.Hostname,
		OS:       h.OS,
		Platform: h.Platform,
		Uptime:   h.Uptime,
	}, nil
}
