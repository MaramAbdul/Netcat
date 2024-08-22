package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn     net.Conn
	username string
}

var (
	clients  = make(map[net.Conn]*Client)
	messages []string
	mutex    = &sync.Mutex{}
)

func main() {
	server, err := net.Listen("tcp", ":8989")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer server.Close()

	fmt.Println("Server started on port 8989")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func readWelcomeMessage() string {
	content, err := ioutil.ReadFile("welcome.txt")
	if err != nil {
		return "Welcome to the chat!\n"
	}
	return string(content)

}

func handleClient(conn net.Conn) {
	defer conn.Close()

	mutex.Lock()
	if len(clients) >= 10 {
		mutex.Unlock()
		conn.Write([]byte("Server is full. Please try again later.\n"))
		return
	}
	mutex.Unlock()

	// here  we will display the welcome message 
	writer := bufio.NewWriter(conn)
	welcomeMessage := readWelcomeMessage()

	lines := bufio.NewScanner(strings.NewReader(welcomeMessage))
	for lines.Scan() {
		writer.WriteString(lines.Text() + "\n")
		writer.Flush() 
	}

	writer.Flush()

	reader := bufio.NewReader(conn)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	mutex.Lock()
	clients[conn] = &Client{conn: conn, username: username}

	notifyAll(fmt.Sprintf("%s joined our chat...\n", username))
	mutex.Unlock()

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		message = strings.TrimSpace(message)
		if message != "" {
			timestampedMessage := fmt.Sprintf("[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), username, message)
			saveAndBroadcastMessage(timestampedMessage, conn)
		}
	}

	mutex.Lock()
	notifyAll(fmt.Sprintf("%s left our chat...\n", username))
	delete(clients, conn)
	mutex.Unlock()
}

func saveAndBroadcastMessage(message string, sender net.Conn) {
	mutex.Lock()
	messages = append(messages, message)
	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(message))
		}
	}
	mutex.Unlock()
}

func notifyAll(message string) {
	for conn := range clients {
		conn.Write([]byte(message))
	}
}

func sendChatHistory(conn net.Conn) {
	for _, message := range messages {
		conn.Write([]byte(message))
	}
}
