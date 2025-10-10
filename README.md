# USDC Event Tracker

A high-performance blockchain event tracker specifically designed to monitor USDC (USD Coin) transactions across multiple Ethereum-compatible networks. The tracker provides real-time monitoring with configurable data sinks for storage and processing.

## Features

- Real-time block monitoring
- USDC transaction filtering
- Event detection (Transfer, Approval)
- Graceful shutdown handling
- Structured logging for better readability

### 📊 Flexible Data Sinks
- **Console** - Real-time console output with formatted display
- **Filesystem** - Multiple formats (JSON, JSONL, CSV, Text) with rotation
- **PostgreSQL** - Structured database storage with indexing
- **MongoDB** - Document-based storage with flexible querying
- **Apache Kafka** - Event streaming with partitioning and compression

### 🚀 Performance Features
- **Batch Processing** - Configurable batch sizes for optimal performance
- **Background Workers** - Async processing with graceful shutdown
- **Connection Pooling** - Efficient resource management
- **Error Handling** - Comprehensive error recovery and logging
- **Metrics** - Built-in performance and health metrics

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Access to an Ethereum RPC endpoint (Infura, Alchemy, or self-hosted)

### Installation

```bash
# Clone the repository
git clone https://github.com/your-username/usdc-event-tracker.git
cd usdc-event-tracker

# Install dependencies
go mod tidy

# Copy environment template
cp .env.example .env
```

### Configuration

Edit `.env` file with your settings:

```env
# Ethereum connection (required)
WEBHOOK_URL=wss://sepolia.infura.io/ws/v3/YOUR_PROJECT_ID

# Network selection
NETWORK=sepolia

# Data sinks (comma-separated)
SINKS=console,filesystem

# Filesystem sink settings
FS_OUTPUT_DIR=./usdc-events
FS_FORMAT=json
FS_FILE_PREFIX=sepolia-usdc
```

### Run

```bash
# Build the application
go build -o usdc-event-tracker .

# Run with your configuration
./usdc-event-tracker
```

## Configuration Reference

### Environment Variables

| Variable | Description | Default | Options |
|----------|-------------|---------|---------|
| `WEBHOOK_URL` | Ethereum RPC endpoint (required) | - | HTTP/HTTPS/WS/WSS URL |
| `NETWORK` | Blockchain network | `sepolia` | `mainnet`, `sepolia`, `arbitrum`, `optimism`, `polygon`, `avalanche`, `linea` |
| `SINKS` | Data output destinations | `console` | `console`, `filesystem`, `sql`, `mongodb`, `kafka` |

### Filesystem Sink

| Variable | Description | Default | Options |
|----------|-------------|---------|---------|
| `FS_OUTPUT_DIR` | Output directory | `./usdc-events` | Any valid path |
| `FS_FORMAT` | File format | `json` | `json`, `jsonl`, `csv`, `text` |
| `FS_FILE_PREFIX` | File name prefix | `usdc-events` | Any string |

### PostgreSQL Sink

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SQL_CONNECTION_STRING` | PostgreSQL connection string | - | ✅ |
| `SQL_TABLE_NAME` | Table name for events | `usdc_events` | ❌ |
| `SQL_SCHEMA_NAME` | Database schema | `public` | ❌ |
| `SQL_BATCH_SIZE` | Batch size for inserts | `100` | ❌ |
| `SQL_CREATE_TABLES` | Auto-create tables | `true` | ❌ |

### MongoDB Sink

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `MONGO_URI` | MongoDB connection URI | - | ✅ |
| `MONGO_DATABASE` | Database name | `blockchain` | ❌ |
| `MONGO_COLLECTION` | Events collection | `usdc_events` | ❌ |
| `MONGO_LOGS_COLLECTION` | Logs collection | `usdc_logs` | ❌ |
| `MONGO_BATCH_SIZE` | Batch size | `50` | ❌ |

### Kafka Sink

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `KAFKA_BROKERS` | Comma-separated broker list | - | ✅ |
| `KAFKA_TOPIC` | Main topic name | `usdc-events` | ❌ |
| `KAFKA_LOGS_TOPIC` | Separate logs topic | - | ❌ |
| `KAFKA_BATCH_SIZE` | Messages per batch | `100` | ❌ |
| `KAFKA_COMPRESSION` | Compression algorithm | `gzip` | `gzip`, `snappy`, `lz4`, `zstd` |

## Architecture

### Core Components

```
┌─────────────────────┐    ┌──────────────────────┐    ┌─────────────────────┐
│                     │    │                      │    │                     │
│   Ethereum Client   │────│   Event Tracker      │────│   Sink Manager      │
│   (WebSocket/HTTP)  │    │   (Block Monitor)    │    │   (Multi-output)    │
│                     │    │                      │    │                     │
└─────────────────────┘    └──────────────────────┘    └─────────────────────┘
                                        │                           │
                                        │                           │
                                        ▼                           ▼
                           ┌──────────────────────┐    ┌─────────────────────┐
                           │                      │    │                     │
                           │   USDC Filter        │    │   Console Sink      │
                           │   (ERC20 Decoder)    │    │   Filesystem Sink   │
                           │                      │    │   PostgreSQL Sink   │
                           └──────────────────────┘    │   MongoDB Sink      │
                                                       │   Kafka Sink        │
                                                       │                     │
                                                       └─────────────────────┘
```

### Data Flow

1. **Block Monitoring** - Continuously polls for new blocks
2. **Transaction Filtering** - Identifies USDC-related transactions  
3. **Event Decoding** - Decodes Transfer and Approval events
4. **Sink Distribution** - Sends events to all configured sinks
5. **Batch Processing** - Optimizes throughput with batching

### Event Structure

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "blockNumber": 19000000,
  "txHash": "0x1234567890abcdef...",
  "status": 1,
  "gasUsed": 65000,
  "logs": [
    {
      "type": "Transfer",
      "address": "0x...",
      "topics": ["0x...", "0x...", "0x..."],
      "data": "0x...",
      "from": "0x...",
      "to": "0x...",
      "value": "1000000000"
    }
  ]
}
```

## Development

### Project Structure

```
usdc-event-tracker/
├── cmd/                    # Application entrypoints
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── erc20/             # ERC20 event definitions
│   ├── sinks/             # Data output implementations
│   │   ├── console/       # Console output
│   │   ├── fs/            # Filesystem output
│   │   ├── sql/           # PostgreSQL output
│   │   ├── mongodb/       # MongoDB output
│   │   └── kafka/         # Kafka output
│   ├── tracker/           # Core tracking logic
│   ├── tx/                # Transaction utilities
│   ├── usdc/              # USDC-specific utilities
│   └── ws/                # WebSocket client
├── .env.example           # Environment template
├── go.mod                 # Go module definition
└── main.go               # Application entry point
```

### Dependencies

The project includes skeleton implementations for all sinks. To implement them, you'll need:

#### Core Dependencies (already included)
```bash
go get github.com/ethereum/go-ethereum
go get github.com/joho/godotenv
```

#### PostgreSQL Sink
```bash
go get github.com/lib/pq
```

#### MongoDB Sink  
```bash
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
```

#### Kafka Sink
```bash
go get github.com/segmentio/kafka-go
```

All sinks implement the `Sink` interface:

```go
type Sink interface {
    Name() string
    Initialize() error
    Write(ctx context.Context, events []Event) error
    Close() error
}
```

### Building

```bash
# Install dependencies for all sinks
go mod tidy

# Build the application
go build -o usdc-event-tracker .

# Run tests
go test ./...
```

## Production Deployment

### Docker Example

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o usdc-event-tracker .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/usdc-event-tracker .
CMD ["./usdc-event-tracker"]
```

### Environment Variables Summary

#### Required
- `WEBHOOK_URL` - Ethereum RPC endpoint

#### Optional
- `NETWORK` - Network to monitor (default: sepolia)
- `SINKS` - Output destinations (default: console)

#### Sink-Specific
- **Filesystem**: `FS_OUTPUT_DIR`, `FS_FORMAT`, `FS_FILE_PREFIX`
- **PostgreSQL**: `SQL_CONNECTION_STRING`, `SQL_TABLE_NAME`, etc.
- **MongoDB**: `MONGO_URI`, `MONGO_DATABASE`, etc.
- **Kafka**: `KAFKA_BROKERS`, `KAFKA_TOPIC`, etc.

## Implementation Progress

### ✅ Completed
- Multi-network support (7 networks)
- ERC20 event decoding (Transfer, Approval)
- Sink architecture with interface
- Console sink (full implementation)
- Filesystem sink (skeleton + full reference)

### 🚧 To Implement (Skeletons Ready)
- **SQL Sink** - PostgreSQL with batch processing
- **MongoDB Sink** - Document storage with indexing  
- **Kafka Sink** - Event streaming with compression

Each skeleton includes:
- Complete struct definitions
- Function signatures with detailed TODO comments
- Import statements and dependencies
- Full implementation reference in `/tmp/`

## Contributing

1. Fork the repository
2. Implement sink skeletons following the TODO comments
3. Test your implementation thoroughly
4. Submit a pull request


## License

The tracker displays:
- Connection information
- Block processing details
- USDC transactions with:
  - Transaction hash
  - Status (success/failure)
  - Gas usage
  - Event types (Transfer/Approval)
  - Event details (from/to addresses)
