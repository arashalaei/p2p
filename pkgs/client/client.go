package client

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"zarban.io/p2p/pkgs/discovery"
	"zarban.io/p2p/pkgs/pubsub"
)

type Node struct {
	host      host.Host
	topicName string
}

type Config struct {
	ListenAddr string
	TopicName  string
}

func NewNode(config Config) (*Node, error) {

	host, err := libp2p.New(
		libp2p.ListenAddrStrings(config.ListenAddr),
	)
	if err != nil {
		return nil, err
	}

	return &Node{
		host:      host,
		topicName: config.TopicName,
	}, nil
}

func (n *Node) Run() {
	ctx := context.Background()
	// Create a discovery with the given host and topic
	d, err := discovery.NewDiscovery(ctx, n.host, n.topicName)
	if err != nil {
		fmt.Printf("Failed to create discovery: %v\n", err)
		return
	}

	// Discover and connect to peers that are subscribed to the topic
	go d.DiscoverPeers()

	// Create a pubsub with the given host and options
	p, err := pubsub.NewPubSub(ctx, n.host)
	if err != nil {
		fmt.Printf("Failed to create pubsub: %v\n", err)
		return
	}

	// Join the topic and get a topic handle
	t, err := p.JoinTopic(n.topicName)
	if err != nil {
		fmt.Printf("Failed to join topic: %v\n", err)
		return
	}

	// Stream the console input to the topic
	go t.SendTo(ctx)

	sub, err := t.Subscribe()
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	// Print the messages from the topic
	t.GetFrom(ctx, sub)
}
