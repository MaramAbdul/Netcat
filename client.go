package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8989")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// here will recive the message from the server
	go receiveMessages(conn)
	//here will send the message to the server
	sendMessages(conn)
}

func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed by server.")
			return
		}
		fmt.Print(message)
	}
}

func sendMessages(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if strings.TrimSpace(message) == "" {
			continue
		}
		conn.Write([]byte(message + "\n"))
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from input:", err)
	}
}
