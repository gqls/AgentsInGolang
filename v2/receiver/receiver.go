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

type ReceiverAgent struct {
	id       string
	messages []Message
	conn     net.Conn
}

func NewReceiverAgent(id string, ln net.Listener) (*ReceiverAgent, error) {
	conn, err := ln.Accept()
	if err != nil {
		return nil, err
	}

	return &ReceiverAgent{
		id:       id,
		messages: make([]Message, 0),
		conn:     conn,
	}, nil
}

func (a *ReceiverAgent) ReceiveMessage() {
	decoder := json.NewDecoder(a.conn)
	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			fmt.Printf("Error receiving message: %v\n", err)
			return
		}
		a.messages = append(a.messages, msg)
		fmt.Printf("ReceiverAgent %s received: '%s' from %s\n", a.id, msg.Content, msg.SenderID)

		// Send response
		response := fmt.Sprintf("Received your message: %s", msg.Content)
		a.SendMessage(response)
	}
}

func (a *ReceiverAgent) SendMessage(content string) error {
	msg := Message{
		SenderID:  a.id,
		Content:   content,
		Timestamp: time.Now(),
	}

	encoder := json.NewEncoder(a.conn)
	return encoder.Encode(msg)
}

func (a *ReceiverAgent) Run() {
	fmt.Printf("ReceiverAgent %s started\n", a.id)
	a.ReceiveMessage()
}

func main() {
	// Set up TCP listener
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	defer ln.Close()

	fmt.Println("Receiver server listening on :8080")

	receiver, err := NewReceiverAgent("Agent2", ln)
	if err != nil {
		fmt.Printf("Failed to create receiver: %v\n", err)
		return
	}

	receiver.Run()
	fmt.Println("Receiver simulation complete")
}
