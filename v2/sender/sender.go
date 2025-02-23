package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Message struct {
	SenderID  string
	Content   string
	Timestamp time.Time
}

type SenderAgent struct {
	id       string
	messages []Message
	conn     net.Conn
}

func NewSenderAgent(id string, serverAddr string) (*SenderAgent, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	return &SenderAgent{
		id:       id,
		messages: make([]Message, 0),
		conn:     conn,
	}, nil
}

func (a *SenderAgent) ReceiveMessage() {
	decoder := json.NewDecoder(a.conn)
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

	encoder := json.NewEncoder(a.conn)
	return encoder.Encode(msg)
}

func (a *SenderAgent) Run() {
	fmt.Printf("SenderAgent %s started\n", a.id)

	// Start receiving messages in a goroutine
	go a.ReceiveMessage()

	// Send messages
	for i := 1; i <= 3; i++ {
		message := fmt.Sprintf("Hello message #%d", i)
		err := a.SendMessage(message)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			return
		}
		time.Sleep(1 * time.Second)
	}

	// Keep running to receive responses
	time.Sleep(5 * time.Second)
	a.conn.Close()
}

func main() {
	// Replace with actual receiver address:port
	sender, err := NewSenderAgent("Agent1", "localhost:8080")
	if err != nil {
		fmt.Printf("Failed to create sender: %v\n", err)
		return
	}

	sender.Run()
	fmt.Println("Sender simulation complete")
}
