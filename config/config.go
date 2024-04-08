package config

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
	"gopkg.in/yaml.v2"
)

type Config struct {
	PrivateKey string `yaml:"private_key"`
}

func GetPrivateKey(path string) (crypto.PrivKey, error) {
	// Open the YAML file
	file, err := os.Open("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to open YAML file: %v", err)
	}
	defer file.Close()

	// Parse the YAML data
	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}

	// Decode the hexadecimal string to bytes
	privateKeyBytes, err := hex.DecodeString(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key string: %v", err)
	}

	// Unmarshal the private key bytes
	privateKey, err := crypto.UnmarshalPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private key: %v", err)
	}

	return privateKey, nil
}
