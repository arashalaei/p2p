# P2P Application

## Description

This is a peer-to-peer (P2P) application built in Go using the libp2p library. It provides a command-line interface (CLI) for running bootstrap nodes and client peers, enabling decentralized networking with peer discovery via Kademlia DHT and message broadcasting via Pub/Sub.

The application is designed to demonstrate core P2P concepts, including secure peer identity management, network bootstrapping, and real-time communication. This project was developed as part of my portfolio to showcase skills in Go programming, distributed systems, and networking protocols.

## Features

- **CLI Interface**: Built with Cobra for easy command management.
- **Peer Discovery**: Uses libp2p's Kademlia DHT for efficient peer finding.
- **Pub/Sub Messaging**: Supports publish-subscribe pattern for broadcasting messages across the network.
- **Secure Identity**: Peer private keys loaded from YAML configuration.
- **Docker Support**: Containerized deployment for easy setup and scalability.
- **Modular Structure**: Organized packages for bootstrap, client, discovery, and pubsub functionalities.

## Technologies Used

- **Go**: Primary programming language (version 1.21.0).
- **libp2p**: Core P2P networking library, including modules for DHT and Pub/Sub.
- **Cobra**: For building the CLI.
- **YAML**: Configuration management.
- **Docker**: For containerization.
- Other dependencies: multiaddr, multiformats, and various cryptography and utility libraries (see `go.mod` for full list).

## Installation

### Prerequisites
- Go 1.21 or higher installed.
- Docker (optional, for containerized runs).

### From Source
1. Clone the repository:
   ```
   git clone https://github.com/arashalaei/p2p.git
   cd p2p
   ```
2. Install dependencies:
   ```
   go mod download
   ```
3. Build the application:
   ```
   go build -o p2p-app main.go
   ```

### Using Docker
1. Build the Docker image:
   ```
   docker build -t p2p-app .
   ```
2. Run the container (example for bootstrap command):
   ```
   docker run p2p-app bootstrap
   ```

## Usage

The application provides a CLI with the following commands:

- **Bootstrap Node**: Runs a bootstrap peer to help other nodes join the network.
  ```
  ./p2p-app bootstrap
  ```

- **Client Peer**: Runs a regular peer that connects to the network and participates in Pub/Sub.
  ```
  ./p2p-app client
  ```

Configuration is loaded from `config.yaml`, which includes the peer's private key. Customize this file as needed for different environments.

For full command options, use:
```
./p2p-app --help
```

## Project Structure

- `cmd/`: CLI commands (root, bootstrap, client).
- `config/`: Configuration loading logic.
- `pkgs/`: Core packages for bootstrap, client, discovery, and pubsub.
- `main.go`: Entry point.
- `config.yaml`: Default configuration file.
- `Dockerfile`: For building Docker images.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes. Ensure code follows Go best practices and includes tests where applicable.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details (note: add a LICENSE file if not present).

## Contact

For questions or collaboration, reach out via GitHub issues or my portfolio site.
