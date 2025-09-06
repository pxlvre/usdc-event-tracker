# USDC Event Tracker

A Go application that monitors Ethereum blockchain for USDC token transactions on the Sepolia testnet.

## Features

- Real-time block monitoring
- USDC transaction filtering
- Event detection (Transfer, Approval)
- Graceful shutdown handling
- Structured logging with emojis for better readability

## Project Structure

```
usdc-event-tracker/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── tracker/           # Main tracking logic
│   │   └── tracker.go
│   ├── tx/               # Transaction handling
│   │   └── tx.go
│   ├── usdc/             # USDC-specific filtering
│   │   └── usdc.go
│   └── ws/               # WebSocket/Ethereum client
│       └── eth_ws.go
├── .env                  # Environment variables
├── go.mod
└── go.sum
```

## Setup

1. Clone the repository
2. Create a `.env` file with your Infura/Alchemy endpoint:
   ```
   WEBHOOK_URL=https://sepolia.infura.io/v3/YOUR_API_KEY
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

## Running

```bash
go run main.go
```

Or build and run:
```bash
go build -o usdc-tracker
./usdc-tracker
```

## Architecture

The application follows Go best practices with:

- **Clean Architecture**: Separation of concerns with internal packages
- **Error Handling**: Proper error wrapping and propagation
- **Context Support**: Graceful shutdown with context cancellation
- **Modular Design**: Each package has a single responsibility

### Key Components

- **Config**: Manages environment variables and application settings
- **Tracker**: Orchestrates the monitoring process
- **TX**: Handles blockchain transaction retrieval
- **USDC**: Filters and processes USDC-specific transactions
- **WS**: Manages Ethereum client connections

## USDC Contract

The tracker monitors the official USDC contract on Sepolia:
- Address: `0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238`

## Output

The tracker displays:
- Connection information
- Block processing details
- USDC transactions with:
  - Transaction hash
  - Status (success/failure)
  - Gas usage
  - Event types (Transfer/Approval)
  - Event details (from/to addresses)