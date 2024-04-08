package bootstrap

import (
	"log"

	"github.com/spf13/cobra"
	"zarban.io/p2p/config"
	"zarban.io/p2p/pkgs/bootstrap"
)

const (
	DefaultListenAddr = "/ip4/0.0.0.0/tcp/8080"
)

func main() {
	privateKey, err := config.GetPrivateKey("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	cfg := bootstrap.Config{
		ListenAddr:  DefaultListenAddr,
		RelayServer: true,
		PrivateKey:  privateKey,
	}
	node, err := bootstrap.NewNode(cfg)
	if err != nil {
		log.Fatal(err)
	}
	node.Run()
}

func Register(root *cobra.Command) {
	root.AddCommand(
		&cobra.Command{
			Use:   "bootstrap",
			Short: "Run the bootstrap mode",
			Run: func(cmd *cobra.Command, args []string) {
				main()
			},
		},
	)
}
