package main

import (
	"bufio"
	"fmt"
	"net"
)

const bufferSize = 5

func main() {

	var ringBuffer [bufferSize]string
	writeIndex := 0
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting network socket:", err)
		return
	}

	defer listener.Close()
	fmt.Println("Ingestion engine initialized. Listening on port 8080...")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connnection:", err)
		return
	}

	defer conn.Close()
	fmt.Println("Client connected! Processing incoming network stream")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		logEntry := scanner.Text()

		ringBuffer[writeIndex] = logEntry
		fmt.Println("[Memory allocation] Stored log at index %d: %s\n", writeIndex, logEntry)

		writeIndex = (writeIndex + 1) % bufferSize

		fmt.Println("\n---- Final In-Memory State of the Ring Buffer---")
		for i, log := range ringBuffer {
			fmt.Printf(":Slot[%d]: %s\n", i, log)
		}
	}
}
