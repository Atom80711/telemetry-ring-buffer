package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const bufferSize = 1000

func main() {

	file, err := os.Create("production.log")
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}

	defer file.Close()

	fileWriter := bufio.NewWriter(file)

	//Create a thread safe buffered channel to hold logs in RAM
	logChannel := make(chan string, bufferSize)

	go backgroundFlusher(logChannel, fileWriter)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting network socket:", err)
		return
	}

	defer listener.Close()

	conn, err := listener.Accept()

	if err != nil {
		return
	}

	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		//straight into channel from network socket without any locks
		logChannel <- scanner.Text()
	}

	close(logChannel)
}
func backgroundFlusher(logChannel chan string, writer *bufio.Writer) {
	// This loop automatically blocks and sleeps when the channel is empty,
	// and wakes up instantly the microsecond a new log enters the channel.
	for logEntry := range logChannel {
		writer.WriteString(logEntry + "\n")
		writer.Flush()
	}
}
