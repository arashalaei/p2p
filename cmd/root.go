package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"zarban.io/p2p/cmd/bootstrap"
	"zarban.io/p2p/cmd/client"
)

func Execute() {
	root := &cobra.Command{
		Use:     "p2p",
		Short:   "p2p is a peer-to-peer application ",
		Version: "0.1",
	}

	bootstrap.Register(root)
	client.Register(root)

	if err := root.Execute(); err != nil {
		log.Panic("error while executing command")
	}
}
