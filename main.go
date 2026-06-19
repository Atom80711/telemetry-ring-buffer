package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const bufferSize = 10

var (
	ringBuffer [bufferSize]string
	writeIndex = 0
	readIndex  = 0
	mutex      sync.Mutex
)

func main() {

	file, err := os.Create("production.log")
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}

	defer file.Close()
	fileWriter := bufio.NewWriter(file)

	go backgroundFlusher(fileWriter)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting network socket", err)
		return
	}

	defer listener.Close()
	fmt.Println("Ingestion Engine running with Async Worker. Listening on port 8080...")

	conn, err := listener.Accept()
	if err != nil {
		return
	}

	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		logEntry := scanner.Text()
		//Secure the lock before writing to shared memory
		mutex.Lock()
		ringBuffer[writeIndex] = logEntry
		writeIndex = (writeIndex + 1) % bufferSize
		mutex.Unlock() //Release instantly so network layer isnt blocked

	}
}

func backgroundFlusher(writer *bufio.Writer) {
	//Wake up every 500 milliseconds to flush memory data to disk hardware
	for {

		time.Sleep(500 * time.Millisecond)

		mutex.Lock()
		//Process all pending unread items currently sitting in the ring buffer

		for readIndex != writeIndex {
			logToFlush := ringBuffer[readIndex]

			//Write directly to our high performance file buffer
			writer.WriteString(logToFlush + "\n")

			readIndex = (readIndex + 1) % bufferSize
		}

		// Ensure data is completely pushed out of RAN buffer
		writer.Flush()
		mutex.Unlock()
	}
}
