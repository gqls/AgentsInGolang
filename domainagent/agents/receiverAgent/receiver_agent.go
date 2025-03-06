package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Message struct {
	SenderID  string
	Content   string
	Timestamp time.Time
}

type ReceiverAgent struct {
	id         string
	messages   []Message
	connection net.Conn
}

func NewReceiverAgent(id string, connection net.Conn) *ReceiverAgent {
	return &ReceiverAgent{
		id:         id,
		messages:   make([]Message, 0),
		connection: connection,
	}
}

func (a *ReceiverAgent) ReceiveMessage() {
	decoder := json.NewDecoder(a.connection)
	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			fmt.Printf("Error receiving message: %v\n", err)
			return
		}
		a.messages = append(a.messages, msg)
		fmt.Printf("ReceiverAgent %s received: '%s' from %s\n", a.id, msg.Content, msg.SenderID)

		// send response
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

	encoder := json.NewEncoder(a.connection)
	return encoder.Encode(msg)
}

func (a *ReceiverAgent) Run() {
	fmt.Printf("RecieverAgent %s started \n", a.id)
	a.ReceiveMessage()
}

func main() {
	// set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// set up TCP listener
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Failed to start receiver listening server on port 8080 %v\n", err)
		return
	}
	defer ln.Close()

	fmt.Println("Receiver server listening on :8080")

	// channel to track active connections
	connections := make(chan net.Conn)

	// Accept connections in go routine
	go func() {
		for {
			connection, err := ln.Accept()
			if err != nil {
				fmt.Printf("Error accepting connection: %v\n", err)
				continue
			}
			connections <- connection
		}
	}()

	// Keep track of active agents
	agents := make([]*ReceiverAgent, 0)

	// Main Loop
	for {
		select {
		case connection := <-connections:
			// New connection received
			agentID := fmt.Sprintf("ReceiverAgetnt-%d", len(agents)+1)
			agent := NewReceiverAgent(agentID, connection)
			agents = append(agents, agent)
			go agent.Run()

		case <-sigs:
			// shutdown signal received
			fmt.Println("Received shutdown signal to Receiver agent. Shutting down...")
			return
		}
	}
}
