# Transaction Parser for Ethereum

This project is a transaction parser that interacts with the Ethereum blockchain, fetching blocks, transactions, and handling subscription and notification services for addresses. The application is built with Go, using an in-memory storage system, and interacts with Ethereum using the JSON-RPC API.

## Features

- **Fetch current block**: Fetches the latest block from the Ethereum blockchain.
- **Fetch block by number**: Retrieves a specific block and its transactions by block number.
- **Subscribe to an address**: Allows users to subscribe to an Ethereum address to track transactions.
- **Track transactions**: Tracks incoming and outgoing transactions for subscribed addresses.
- **In-memory storage**: Stores address subscriptions and transactions using in-memory storage.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Endpoints](#endpoints)
- [Testing](#testing)
- [License](#license)

## Installation

### Prerequisites

- **Go**: Ensure you have Go installed (version 1.18 or higher).
- **Git**: Ensure Git is installed to clone the repository.

### Steps

1. Clone the repository:

```bash
git clone https://github.com/your-repository/tx-parser.git
cd tx-parser
```

2. Install the Go dependencies:

```bash
go mod tidy
```

### Usage

To start the application, you can run the following command:

```bash
go run cmd/parser/main.go
```

By default, the application will start an HTTP server on localhost:8088 and expose several API endpoints for subscribing to addresses and fetching transactions.

### Configuration

The application configuration is stored in the config.yaml file. You can modify it to change the Ethereum RPC URL, logging level, or the server port.

```yaml
server:
   port: ":8088"  # Server port
   host: "localhost"
   ethrpc: "https://ethereum-rpc.publicnode.com"

logging:
   level: "debug"  # Available options: debug, info, warn, error
```

### Project Structure

```bash
.
├── cmd
│   └── parser           # Main application entry point
├── configs              # Configuration files
├── internal
│   ├── api              # HTTP server and route handlers
│   ├── app              # Application setup and main logic
│   ├── config           # Configuration handling
│   ├── interfaces       # Interfaces for parser and storage
│   ├── parser           # Ethereum parser (fetching transactions and blocks)
│   ├── rpc              # Ethereum JSON-RPC client
│   └── storage          # In-memory storage for addresses and transactions
├── pkg
│   └── logger           # Custom logger package
├── scripts              # Any custom scripts
├── utils                # Utility functions (e.g., address normalization)
```

### Endpoints

1. Fetch Current Block
Method: GET
Endpoint: /current-block
Description: Fetches the latest block from the Ethereum blockchain.

Example:
```bash
curl http://localhost:8088/current-block
```

2. Subscribe to an Address
Method: POST
Endpoint: /subscribe
Description: Subscribes to an Ethereum address to track its transactions.
Example:
```bash
curl -X POST http://localhost:8088/subscribe -d '{"address": "0xYourAddress"}' -H 'Content-Type: application/json'
```

3. Get Transactions for an Address
Method: GET
Endpoint: /transactions/{address}
Description: Fetches transactions (incoming and outgoing) for the specified Ethereum address.
Example:
```bash
curl http://localhost:8088/transactions/0xYourAddress
```

### Testing

The project includes unit tests for the core components such as the Ethereum parser, RPC client, and in-memory storage. To run the tests, use the following command:

```bash
go test ./...
```

This will execute all the tests in the project and provide coverage for critical functionality like subscribing to addresses, fetching blocks, and tracking transactions.

### License
This project is licensed under the MIT License. See the LICENSE file for more information.