package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Message struct {
	SenderID  string
	Content   string
	Timestamp time.Time
}

type SenderAgent struct {
	id         string
	messages   []Message
	connection net.Conn
}

func NewSenderAgent(id string, serverAddr string) (*SenderAgent, error) {
	// attempt to connect with retries
	var connection net.Conn
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		connection, err = net.Dial("tcp", serverAddr)
		if err == nil {
			break
		}
		fmt.Printf("Connection attempt %d failed: %v\n", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect after %d attempts: %v", maxRetries, err)
	}

	return &SenderAgent{
		id:         id,
		messages:   make([]Message, 0),
		connection: connection,
	}, nil
}

func (a *SenderAgent) ReceiveMessage() {
	decoder := json.NewDecoder(a.connection)
	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			fmt.Printf("Error receiving message: %v\n", err)
			return
		}
		a.messages = append(a.messages, msg)
		fmt.Printf("SenderAgent %s received: '%s' from %s\n", a.id, msg.Content, msg.SenderID)
	}
}

func (a *SenderAgent) SendMessage(content string) error {
	msg := Message{
		SenderID:  a.id,
		Content:   content,
		Timestamp: time.Now(),
	}

	encoder := json.NewEncoder(a.connection)
	return encoder.Encode(msg)
}

func (a *SenderAgent) Run() {
	fmt.Printf("SenderAgent %s started\n", a.id)

	// start receiving messages in a goroutine
	go a.ReceiveMessage()

	// send periodic messages
	for {
		message := fmt.Sprintf("Hellow from AgentA %s at %s", a.id, time.Now().Format(time.RFC3339))
		err := a.SendMessage(message)
		if err != nil {
			fmt.Printf("AgentA error sending message: %v\n", err)
			return
		}

		time.Sleep(10 * time.Second)
	}
}

func main() {
	// server address from env variable or default
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "127.0.0.1:8080"
	}

	fmt.Printf("AgentA sender agent connecting to server %s\n", serverAddr)

	// Create and Run Sender agent
	sender, err := NewSenderAgent("SenderAgent-1", serverAddr)
	if err != nil {
		fmt.Printf("Failed to create sender agent A: %v\n", err)
		return
	}

	sender.Run()
}
