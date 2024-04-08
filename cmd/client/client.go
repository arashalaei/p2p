package client

import (
	"log"

	"github.com/spf13/cobra"
	"zarban.io/p2p/pkgs/client"
)

func main(listenAddr, topicName string) {
	cfg := client.Config{
		ListenAddr: listenAddr,
		TopicName:  topicName,
	}
	node, err := client.NewNode(cfg)
	if err != nil {
		log.Fatal(err)
	}
	node.Run()
}

func Register(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "Run the client mode",
		Run: func(cmd *cobra.Command, args []string) {
			listenAddr, err := cmd.Flags().GetString("listen-addr")
			if err != nil {
				log.Fatal(err)
			}
			topicName, err := cmd.Flags().GetString("topic")
			if err != nil {
				log.Fatal(err)
			}
			main(listenAddr, topicName)
		},
	}

	cmd.Flags().StringP("listen-addr", "l", "", "Listen address for the client")
	cmd.Flags().StringP("topic", "t", "", "Topic name to subscribe to")
	cmd.MarkFlagRequired("listen-addr")
	cmd.MarkFlagRequired("topic")

	root.AddCommand(cmd)
}
