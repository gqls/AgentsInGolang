//package main

import (
	"fmt"
	"sync"
)

// Agent structure
type Agent struct {
	ID      int
	Task    string
	mailbox chan Message
	quit    chan bool
}

// Message structure
type Message struct {
	SenderID   int
	ReceiverID int
	Content    string
}

// Message broker
type Broker struct {
	agents     map[int]Agent
	messages   chan Message
	brokerQuit chan bool
	agentQuit  map[int]chan bool
	mu         sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		agents:     make(map[int]Agent),
		messages:   make(chan Message, 10), // Buffered channel
		brokerQuit: make(chan bool),
		agentQuit:  make(map[int]chan bool),
	}
}

func (b *Broker) Register(agent Agent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.agents[agent.ID] = agent
	b.agentQuit[agent.ID] = make(chan bool)
}

func (b *Broker) RouteMessages() {
	for {
		select {
		case msg := <-b.messages:
			b.mu.Lock()
			receiver, ok := b.agents[msg.ReceiverID]
			_, receiverChannelOpen := b.agentQuit[msg.ReceiverID]
			b.mu.Unlock()
			if ok && receiverChannelOpen {
				receiver.mailbox <- msg
			} else {
				fmt.Println("Message to unknown agent:", msg.ReceiverID)
			}
		case <-b.brokerQuit:
			fmt.Println("Broker shutting down...")
			b.mu.Lock()
			for _, agent := range b.agents {
				close(agent.mailbox)
				agent.quit <- true
			}
			b.mu.Unlock()
			return
		}
	}
}

func (b *Broker) SendMessage(msg Message) {
	b.messages <- msg
}

// Agent's main loop
func (a *Agent) Run(broker *Broker, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(a.quit)

	messagesToSend := []string{"hello", "how are you?", "goodbye"}
	for _, msgContent := range messagesToSend {
		msg := Message{SenderID: a.ID, ReceiverID: 2, Content: msgContent}
		broker.SendMessage(msg)
	}

	for {
		select {
		case msg := <-a.mailbox:
			fmt.Printf("Agent %d received message from Agent %d: %s\n", a.ID, msg.SenderID, msg.Content)
			switch msg.Content {
			case "hello":
				reply := Message{SenderID: a.ID, ReceiverID: msg.SenderID, Content: "hi"}
				broker.SendMessage(reply)
			case "how are you?":
				reply := Message{SenderID: a.ID, ReceiverID: msg.SenderID, Content: "I'm well thank you"}
				broker.SendMessage(reply)
			case "goodbye":
				reply := Message{SenderID: a.ID, ReceiverID: msg.SenderID, Content: "goodbye"}
				broker.SendMessage(reply)
				fmt.Printf("Agent %d exiting due to 'goodbye' message\n", a.ID)
				return
			default:
				fmt.Printf("Agent %d received unknown message %s\n", a.ID, msg.Content)
			}
		case <-a.quit:
			fmt.Printf("Agent %d shutting down...\n", a.ID)
			return
		}
	}
}

func main() {
	broker := NewBroker()
	var wg sync.WaitGroup

	agent1 := Agent{ID: 1, Task: "Task 1", mailbox: make(chan Message, 10), quit: make(chan bool)}
	agent2 := Agent{ID: 2, Task: "Task 2", mailbox: make(chan Message, 10), quit: make(chan bool)}

	broker.Register(agent1)
	broker.Register(agent2)

	wg.Add(2)
	go agent1.Run(broker, &wg)
	go agent2.Run(broker, &wg)
	go broker.RouteMessages()

	wg.Wait()
	broker.brokerQuit <- true
	fmt.Println("Agents finished.")
}
