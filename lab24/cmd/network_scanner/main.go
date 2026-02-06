package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type DeviceInfo struct {
	Name     string
	Type     string
	Addr     string
	LastSeen time.Time
}

func main() {
	beaconPort := flag.Int("beacon-port", 9999, "UDP port for beaconing")
	flag.Parse()

	log.Printf("Starting Enterprise Device Manager...")

	devices := make(map[string]*DeviceInfo)
	var mu sync.Mutex

	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", *beaconPort))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("UDP Listen Error: %v", err)
	}
	defer conn.Close()

	// 1. Discovery Loop
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, srcAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				continue
			}

			parts := strings.Split(string(buffer[:n]), "|")
			if len(parts) < 3 {
				continue
			}

			name, devType, port := parts[0], parts[1], parts[2]
			serviceAddr := fmt.Sprintf("tcp://%s:%s", srcAddr.IP, port)

			mu.Lock()
			if _, exists := devices[name]; !exists {
				log.Printf("[NEW DEVICE] Found %s (%s) at %s", name, devType, serviceAddr)
			}
			devices[name] = &DeviceInfo{
				Name:     name,
				Type:     devType,
				Addr:     serviceAddr,
				LastSeen: time.Now(),
			}
			mu.Unlock()
		}
	}()

	// 2. Management Loop: Periodically print inventory and ping devices
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		mu.Lock()
		fmt.Printf("\n=== Managed Device Inventory (%d devices) ===\n", len(devices))
		for name, info := range devices {
			if time.Since(info.LastSeen) > 15*time.Second {
				log.Printf("[DISCONNECT] Device %s timed out", name)
				delete(devices, name)
				continue
			}
			fmt.Printf("  [%-20s] Type: %-15s | Addr: %s\n", name, info.Type, info.Addr)
			// go pingDevice(info.Addr) // In a real system, we'd health-check here
		}
		fmt.Println("==============================================")
		mu.Unlock()
	}
}
