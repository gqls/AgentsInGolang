package main

import (
	"fmt"
	"time"
)

// Message struct to hold communication data
type Message struct {
	SenderID  string
	Content   string
	Timestamp time.Time
}

// Agent interface defining basic agent behavior
type Agent interface {
	GetID() string
	ReceiveMessage(msg Message)
	SendMessage(content string, receiver Agent)
	Run()
}

// BaseAgent struct with common properties
type BaseAgent struct {
	id       string
	messages []Message
}

// GetID returns agent's identifier
func (a *BaseAgent) GetID() string {
	return a.id
}

// SenderAgent struct
type SenderAgent struct {
	BaseAgent
	receiver Agent
}

// NewSenderAgent creates a new sender agent
func NewSenderAgent(id string, receiver Agent) *SenderAgent {
	return &SenderAgent{
		BaseAgent: BaseAgent{
			id:       id,
			messages: make([]Message, 0),
		},
		receiver: receiver,
	}
}

// ReceiveMessage handles incoming messages
func (a *SenderAgent) ReceiveMessage(msg Message) {
	a.messages = append(a.messages, msg)
	fmt.Printf("SenderAgent %s received: '%s' from %s\n", a.id, msg.Content, msg.SenderID)
}

// SendMessage sends a message to the receiver
func (a *SenderAgent) SendMessage(content string, receiver Agent) {
	msg := Message{
		SenderID:  a.id,
		Content:   content,
		Timestamp: time.Now(),
	}
	receiver.ReceiveMessage(msg)
}

// Run starts the sender agent's behavior
func (a *SenderAgent) Run() {
	fmt.Printf("SenderAgent %s started\n", a.id)
	for i := 1; i <= 3; i++ {
		message := fmt.Sprintf("Hello message #%d", i)
		a.SendMessage(message, a.receiver)
		time.Sleep(1 * time.Second)
	}
}

// ReceiverAgent struct
type ReceiverAgent struct {
	BaseAgent
	sender Agent
}

// NewReceiverAgent creates a new receiver agent
func NewReceiverAgent(id string, sender Agent) *ReceiverAgent {
	return &ReceiverAgent{
		BaseAgent: BaseAgent{
			id:       id,
			messages: make([]Message, 0),
		},
		sender: sender,
	}
}

// ReceiveMessage handles incoming messages
func (a *ReceiverAgent) ReceiveMessage(msg Message) {
	a.messages = append(a.messages, msg)
	fmt.Printf("ReceiverAgent %s received: '%s' from %s\n", a.id, msg.Content, msg.SenderID)

	// Respond to sender
	response := fmt.Sprintf("Received your message: %s", msg.Content)
	a.SendMessage(response, a.sender)
}

// SendMessage sends a message to the sender
func (a *ReceiverAgent) SendMessage(content string, receiver Agent) {
	msg := Message{
		SenderID:  a.id,
		Content:   content,
		Timestamp: time.Now(),
	}
	receiver.ReceiveMessage(msg)
}

// Run starts the receiver agent's behavior
func (a *ReceiverAgent) Run() {
	fmt.Printf("ReceiverAgent %s started\n", a.id)
	// Receiver just waits for messages in this example
}

func main() {
	// Create agents
	sender := NewSenderAgent("Agent1", nil)     // Temporary nil receiver
	receiver := NewReceiverAgent("Agent2", nil) // Temporary nil sender

	// Set the references to each other
	sender.receiver = receiver
	receiver.sender = sender

	// Run agents in goroutines
	go sender.Run()
	go receiver.Run()

	// Keep main running to observe the interaction
	time.Sleep(5 * time.Second)
	fmt.Println("Simulation complete")
}
