package main

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Topic     string
	Body      string
	Timestamp time.Time
	MessageID string
}

func NewMessage(topic, body string) *Message {
	return &Message{
		Topic: topic,
		Body:  body,
	}
}

func (m *Message) GetTopic() string { return m.Topic }
func (m *Message) GetBody() string  { return m.Body }

// ----------------------------------------------------------
// Subscriber represents a subscriber in the pub-sub system.

type Subscriber struct {
	SubscriberID string          // ID of the subscriber
	Messages     chan *Message   // Messages channel
	Topics       map[string]bool // topics the subscriber is subscribed to
	Active       bool            // is the subscriber active
	mu           sync.Mutex      // mutex for concurrent access
}

var subCounter int
var subCounterMutex sync.Mutex

func GenerateSubscriberID() string {
	subCounterMutex.Lock()
	defer subCounterMutex.Unlock()
	subCounter++
	return fmt.Sprintf("sub-%d", subCounter)
}

func NewSubscriber() (string, *Subscriber) {
	id := GenerateSubscriberID()
	return id, &Subscriber{
		SubscriberID: id,
		Messages:     make(chan *Message, 100), // Fixed buffer size
		Topics:       make(map[string]bool),
		Active:       true,
	}
}

func (s *Subscriber) AddTopic(topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Topics[topic] = true
}

func (s *Subscriber) RemoveTopic(topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Topics, topic)
}

func (s *Subscriber) GetTopics() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	topics := []string{}
	for topic := range s.Topics {
		topics = append(topics, topic)
	}
	return topics
}

func (s *Subscriber) Destruct() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Active {
		close(s.Messages)
		s.Active = false
	}
}

func (s *Subscriber) Signal(msg *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Active {
		select {
		case s.Messages <- msg:
		default:
			fmt.Printf("Subscriber %s: channel full, dropping message\n", s.SubscriberID)
		}
	}
}

func (s *Subscriber) Listen() {
	for msg := range s.Messages {
		// Process the message
		fmt.Printf("Subscriber %s received message on topic: %s, body: %s\n",
			s.SubscriberID, msg.GetTopic(), msg.GetBody())
	}
	fmt.Printf("Subscriber %s: listener stopped\n", s.SubscriberID)
}

// ------------------------------------------------
// Broker represents the pub-sub broker.

type Subscribers map[string]*Subscriber

type Broker struct {
	Subscribers Subscribers            // map of subscriber ID to Subscriber
	topics      map[string]Subscribers // map of topic to subscribers
	mu          sync.RWMutex           // mutex for concurrent access
}

func NewBroker() *Broker {
	return &Broker{
		Subscribers: make(Subscribers),
		topics:      make(map[string]Subscribers),
	}
}

func (b *Broker) AddSubscriber() *Subscriber {
	b.mu.Lock()
	defer b.mu.Unlock()
	id, sub := NewSubscriber()
	b.Subscribers[id] = sub

	// Start listening to messages for this subscriber
	go sub.Listen()

	return sub
}

func (b *Broker) RemoveSubscriber(sub *Subscriber) {
	topics := sub.GetTopics()

	for _, topic := range topics {
		b.Unsubscribe(topic, sub)
	}
	sub.Destruct()
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.Subscribers, sub.SubscriberID)
	fmt.Printf("remove subscriber: %s\n", sub.SubscriberID)
}

func (b *Broker) Broadcast(msg string, topics []string) {
	for _, topic := range topics {
		b.Publish(topic, msg)
	}
}

func (b *Broker) GetSubscribers(topic string) Subscribers {
	b.mu.Lock()
	defer b.mu.Unlock()

	if subscribers, exists := b.topics[topic]; exists {
		// Return a copy to avoid concurrent modification issues
		result := make(Subscribers)
		for id, sub := range subscribers {
			result[id] = sub
		}
		return result
	}

	return make(Subscribers)
}

func (b *Broker) Subscribe(topic string, sub *Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// Initialize topic if it doesn't exist
	if _, exists := b.topics[topic]; !exists {
		b.topics[topic] = make(Subscribers)
	}

	// Add subscriber to topic
	b.topics[topic][sub.SubscriberID] = sub

	// Add topic to subscriber's topic list
	sub.AddTopic(topic)

	fmt.Printf("Subscriber %s subscribed to topic: %s\n", sub.SubscriberID, topic)
}

func (b *Broker) Unsubscribe(topic string, sub *Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if subscribers, exists := b.topics[topic]; exists {
		delete(subscribers, sub.SubscriberID)
		sub.RemoveTopic(topic)
		fmt.Printf("Subscriber %s unsubscribed from topic: %s\n", sub.SubscriberID, topic)

		// Clean up empty topics
		if len(subscribers) == 0 {
			delete(b.topics, topic)
		}
	}
}

func (b *Broker) Publish(topic string, msg string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	message := NewMessage(topic, msg)

	if subscribers, exists := b.topics[topic]; exists {
		fmt.Printf("Publishing to topic '%s': %s (to %d subscribers)\n",
			topic, msg, len(subscribers))

		for _, sub := range subscribers {
			go func(s *Subscriber) {
				s.Signal(message)
			}(sub)
		}
	} else {
		fmt.Printf("Topic '%s' has no subscribers, message dropped: %s\n", topic, msg)
	}
}

func (b *Broker) GetTopics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topics := make([]string, 0, len(b.topics))
	for topic := range b.topics {
		topics = append(topics, topic)
	}
	return topics
}

func (b *Broker) GetSubscriberCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.Subscribers)
}

// ----------------------------------------------------------

func main() {
	fmt.Println("=== Starting Pub-Sub System ===")

	// Create broker
	broker := NewBroker()

	// Create subscribers
	sub1 := broker.AddSubscriber()
	sub2 := broker.AddSubscriber()
	sub3 := broker.AddSubscriber()

	fmt.Printf("\nCreated %d subscribers\n", broker.GetSubscriberCount())

	// Subscribe to topics
	fmt.Println("\n=== Subscribing to Topics ===")
	broker.Subscribe("news", sub1)
	broker.Subscribe("news", sub2)
	broker.Subscribe("sports", sub2)
	broker.Subscribe("sports", sub3)
	broker.Subscribe("tech", sub1)
	broker.Subscribe("tech", sub3)

	// Publish messages
	fmt.Println("\n=== Publishing Messages ===")

	// Single topic publishing
	broker.Publish("news", "Breaking: New AI Model Released!")
	time.Sleep(100 * time.Millisecond)

	broker.Publish("sports", "Football: Team A wins championship!")
	time.Sleep(100 * time.Millisecond)

	broker.Publish("tech", "Go 1.21 released with new features!")
	time.Sleep(100 * time.Millisecond)

	// Broadcast to multiple topics
	fmt.Println("\n=== Broadcasting to Multiple Topics ===")
	broker.Broadcast("URGENT: System maintenance scheduled", []string{"news", "tech"})
	time.Sleep(100 * time.Millisecond)

	// Demonstrate unsubscribe
	fmt.Println("\n=== Unsubscribing ===")
	broker.Unsubscribe("news", sub2)
	broker.Publish("news", "This message won't reach sub-2")
	time.Sleep(100 * time.Millisecond)

	// Show current state
	fmt.Println("\n=== System Status ===")
	fmt.Printf("Active topics: %v\n", broker.GetTopics())
	fmt.Printf("Active subscribers: %d\n", broker.GetSubscriberCount())

	// Get subscribers for a topic
	newsSubs := broker.GetSubscribers("news")
	fmt.Printf("Subscribers to 'news': %d\n", len(newsSubs))
	for id := range newsSubs {
		fmt.Printf("  - %s\n", id)
	}

	// Remove a subscriber
	fmt.Println("\n=== Removing Subscriber ===")
	broker.RemoveSubscriber(sub3)
	fmt.Printf("Remaining subscribers: %d\n", broker.GetSubscriberCount())

	// Publish after removal
	broker.Publish("sports", "Tennis: New tournament announced")
	time.Sleep(100 * time.Millisecond)

	// Keep system running to process all messages
	fmt.Println("\n=== Waiting for all messages to be processed ===")
	time.Sleep(2 * time.Second)

	fmt.Println("\n=== Pub-Sub System Demo Complete ===")
}
