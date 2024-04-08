// Package discovery provides a wrapper around the libp2p routing discovery interface.
package discovery

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	ma "github.com/multiformats/go-multiaddr"
)

// Discovery is a type that implements the discovery interface and provides some additional methods.
type Discovery struct {
	ctx              context.Context
	host             host.Host
	routingDiscovery drouting.RoutingDiscovery
	topic            string
}

// NewDiscovery creates a new discovery with the given host and topic.
func NewDiscovery(ctx context.Context, host host.Host, topic string) (*Discovery, error) {
	// Define custom bootstrap nodes
	bootstrapPeers := []peer.AddrInfo{
		{
			ID: "QmQiNTcP9yLAhgMSeh7hf524SRDKHdE3pu8jWT3Ez8xvsY",
			Addrs: []ma.Multiaddr{
				ma.StringCast("/ip4/127.0.0.1/tcp/8080/p2p/QmQiNTcP9yLAhgMSeh7hf524SRDKHdE3pu8jWT3Ez8xvsY"),
			}},
		// Add more bootstrap nodes if needed
	}

	// Create DHT options with custom bootstrap nodes
	dhtOptions := []dht.Option{
		dht.BootstrapPeers(bootstrapPeers...),
	}

	// Initialize a DHT instance
	kademliaDHT, err := dht.New(ctx, host, dhtOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create DHT: %w", err)
	}

	// Bootstrap the DHT with the default bootstrap peers
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, fmt.Errorf("failed to bootstrap DHT: %w", err)
	}

	// Connect to the bootstrap peers in parallel
	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr.Addrs[0])
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Printf("Bootstrap warning: failed to connect to %s, error: %s\n", peerinfo.ID, err)
			}
		}()
	}
	wg.Wait()

	// Create a routing discovery instance with the DHT
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, topic)

	// Wrap the routing discovery in a Discovery type and return it
	return &Discovery{
		ctx:              ctx,
		host:             host,
		routingDiscovery: *routingDiscovery,
		topic:            topic,
	}, nil

}

// FindPeers returns a channel of peer addresses that are subscribed to the topic.
func (d *Discovery) FindPeers() (<-chan peer.AddrInfo, error) {
	// Use the underlying routing discovery to find peers
	return d.routingDiscovery.FindPeers(d.ctx, d.topic)
}

func (d *Discovery) DiscoverPeers() {
	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		log.Println("Searching for peers...")
		peerChan, err := d.routingDiscovery.FindPeers(d.ctx, d.topic)
		if err != nil {
			log.Printf("Failed to find peers: %s\n", err)
			return
		}
		for peer := range peerChan {
			if peer.ID == d.host.ID() {
				continue // No self connection
			}
			err := d.host.Connect(d.ctx, peer)
			if err != nil {
				log.Printf("Failed connecting to %s, error: %s\n", peer.ID, err)
			} else {
				log.Println("Connected to:", peer.ID)
				anyConnected = true
			}
		}
		// Sleep for a while before retrying (optional)
		time.Sleep(5 * time.Second)
	}
	log.Println("Peer discovery complete")
}

// TODO: Add any other methods you need for the discovery type
