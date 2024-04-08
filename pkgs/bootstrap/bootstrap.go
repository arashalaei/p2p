// bootstrap/node.go
package bootstrap

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

const ProtocolVersion = "/bootstrap/1.0.0"

type Node struct {
	host        host.Host
	peerCount   int
	mutex       sync.Mutex
	peerStore   peerstore.Peerstore
	relayServer bool
}

type Config struct {
	ListenAddr  string
	RelayServer bool
	PrivateKey  crypto.PrivKey
}

func NewNode(config Config) (*Node, error) {

	// Create a new libp2p host with relay capabilities
	host, err := libp2p.New(
		libp2p.ListenAddrStrings(config.ListenAddr),
		libp2p.Identity(config.PrivateKey),
		libp2p.EnableRelay(),
	)
	if err != nil {
		return nil, err
	}

	return &Node{
		host:        host,
		peerStore:   host.Peerstore(),
		relayServer: config.RelayServer,
	}, nil
}

func (n *Node) Run() {
	// Set up a stream handler for incoming connections
	n.host.SetStreamHandler(ProtocolVersion, n.handleStream)

	// Set up peer connected and disconnected event handlers
	n.host.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(net network.Network, conn network.Conn) {
			remotePeerID := conn.RemotePeer()
			n.incrementPeerCount()
			log.Printf("Client connected: %s", remotePeerID.String())
			log.Printf("Total connected peers: %d", n.getPeerCount())
		},
		DisconnectedF: func(net network.Network, conn network.Conn) {
			remotePeerID := conn.RemotePeer()
			n.decrementPeerCount()
			log.Printf("Client disconnected: %s", remotePeerID.String())
			log.Printf("Total connected peers: %d", n.getPeerCount())
		},
	})

	// Print the bootstrap node's listening addresses
	fmt.Println("Bootstrap node listening on:")
	for _, addr := range n.host.Addrs() {
		fmt.Printf("%s/p2p/%s\n", addr.String(), n.host.ID().String())
	}

	// If the node is configured as a relay server, announce it as a relay
	if n.relayServer {
		relayAddr, err := multiaddr.NewMultiaddr("/p2p-circuit")
		if err != nil {
			log.Fatal(err)
		}
		relayInfo := peer.AddrInfo{
			ID:    n.host.ID(),
			Addrs: []multiaddr.Multiaddr{relayAddr},
		}
		n.peerStore.AddAddrs(relayInfo.ID, relayInfo.Addrs, 24*time.Hour)
		log.Printf("Announced as relay server: %s", relayInfo.ID.String())
	}

	// Block indefinitely
	select {}
}

func (n *Node) handleStream(stream network.Stream) {
	// Get the remote peer's ID
	remotePeerID := stream.Conn().RemotePeer()

	// Increment the peer count
	log.Printf("New client connected: %s", remotePeerID.String())
	log.Printf("Total connected peers: %d", n.getPeerCount())

	// Handle the stream as needed
	// ...

	// Close the stream when done
	stream.Close()
}

func (n *Node) incrementPeerCount() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.peerCount++
}

func (n *Node) decrementPeerCount() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.peerCount--
}

func (n *Node) getPeerCount() int {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return n.peerCount
}
