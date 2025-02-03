package main

import (
	"fmt"
	"github.com/rebay1982/hflogger/internal/server"
	"encoding/json"
)

func main() {

	server, _ := server.NewWSJTXServer("127.0.0.1", 2237)
	defer server.Close()

	fmt.Println("Looking for WSJT-X...")
	for {
		message, err := server.ReadFromUDP()
		if err != nil {
			fmt.Printf("Failed to read from UDP: %v\n", err)
			continue
		}

		if message.Header.MsgType.String() == "Heartbeat" {
			fmt.Println("Found WSJT-X")
			break
		}
	}

	for {
		message, err := server.ReadFromUDP()
		if err != nil {
		 fmt.Printf("Failed to read from UDP: %v\n", err)
		 continue
		}

		msgJSON, _ := json.Marshal(message)

		fmt.Printf("Got message of type [%s]\n%s\n", message.Header.MsgType, msgJSON)
	}
}
