// Package pubsub provides a wrapper around the libp2p pubsub interface.
package pubsub

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

// Message represents a message containing a pair and price.
type Message struct {
	Pair  string `json:"pair"`
	Price string `json:"price"`
}

// PubSub is a type that implements the pubsub interface and provides some additional methods.
type PubSub struct {
	*pubsub.PubSub
}

// NewPubSub creates a new pubsub with the given host and options.
func NewPubSub(ctx context.Context, h host.Host, opts ...pubsub.Option) (*PubSub, error) {
	// Create a libp2p pubsub instance with the options
	ps, err := pubsub.NewGossipSub(ctx, h, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p pubsub: %w", err)
	}

	// Wrap the libp2p pubsub in a PubSub type and return it
	return &PubSub{ps}, nil
}

// Topic wraps a libp2p pubsub topic.
type Topic struct {
	*pubsub.Topic
}

// JoinTopic joins the given topic and returns a Topic type.
func (p *PubSub) JoinTopic(topic string, opts ...pubsub.TopicOpt) (*Topic, error) {
	// Use the underlying pubsub to join the topic
	t, err := p.PubSub.Join(topic, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to join topic: %w", err)
	}

	// Wrap the libp2p topic in a Topic type and return it
	return &Topic{t}, nil
}

// SendTo sends messages to the topic.
func (t *Topic) SendTo(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter pair: ")
		pair, _ := reader.ReadString('\n')
		fmt.Print("Enter price: ")
		price, _ := reader.ReadString('\n')

		msg := Message{
			Pair:  pair[:len(pair)-1], // Remove the newline character
			Price: price[:len(price)-1],
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("### Marshal error:", err)
			continue
		}

		if err := t.Publish(ctx, msgBytes); err != nil {
			fmt.Println("### Publish error:", err)
		}
	}
}

// GetFrom receives messages from the topic.
func (t *Topic) GetFrom(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("### Next error:", err)
			return
		}

		var msg Message
		if err := json.Unmarshal(m.Message.Data, &msg); err != nil {
			fmt.Println("### Unmarshal error:", err)
			continue
		}

		fmt.Println(m.ReceivedFrom)
		fmt.Printf("%s: %s\n", msg.Pair, msg.Price)
	}
}
