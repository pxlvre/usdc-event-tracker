// Package mongodb implements a MongoDB sink for blockchain events
package mongodb

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ethereum/go-ethereum/core/types"
	"usdc-event-tracker/internal/sinks"
)

// Config holds MongoDB sink configuration
type Config struct {
	URI            string        // MongoDB connection URI
	Database       string        // Database name
	Collection     string        // Collection name for events
	LogsCollection string        // Collection name for event logs
	BatchSize      int           // Number of events to batch before insert
	FlushInterval  time.Duration // Maximum time to wait before flushing batch
	CreateIndexes  bool          // Whether to create indexes
}

// MongoSink writes events to MongoDB
type MongoSink struct {
	config Config
	client *mongo.Client
	
	// Collections
	eventsCollection *mongo.Collection
	logsCollection   *mongo.Collection
	
	// Batch processing
	eventBatch []EventDocument
	logsBatch  []LogDocument
	batchMutex sync.Mutex
	lastFlush  time.Time
	
	// Background processing
	done chan struct{}
	wg   sync.WaitGroup
	
	// Metrics
	totalEvents int64
	totalLogs   int64
	totalBatches int64
}

// EventDocument represents an event in MongoDB
type EventDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Timestamp   time.Time          `bson:"timestamp"`
	BlockNumber uint64             `bson:"blockNumber"`
	TxHash      string             `bson:"txHash"`
	TxStatus    uint64             `bson:"txStatus"`
	GasUsed     uint64             `bson:"gasUsed"`
	EventCount  int                `bson:"eventCount"`
	CreatedAt   time.Time          `bson:"createdAt"`
}

// LogDocument represents an event log in MongoDB
type LogDocument struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	EventID         primitive.ObjectID `bson:"eventId,omitempty"`
	BlockNumber     uint64             `bson:"blockNumber"`
	TxHash          string             `bson:"txHash"`
	LogIndex        uint               `bson:"logIndex"`
	EventType       string             `bson:"eventType"`
	ContractAddress string             `bson:"contractAddress"`
	Topics          []string           `bson:"topics"`
	Data            string             `bson:"data"`
	DecodedData     bson.M             `bson:"decodedData,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt"`
}

// New creates a new MongoDB sink with the given configuration
func New(config Config) *MongoSink {
	// TODO: Implement
	// - Set defaults for empty config values
	// - Initialize MongoSink struct
	// - Setup batch slices and channels
	return nil
}

// Name returns "mongodb" as the sink identifier
func (m *MongoSink) Name() string {
	// TODO: Implement
	return ""
}

// Initialize prepares the MongoDB sink
func (m *MongoSink) Initialize() error {
	// TODO: Implement
	// - Parse MongoDB URI and set connection options
	// - Connect to MongoDB with timeout
	// - Test connection with Ping
	// - Get database and collections
	// - Create indexes if config.CreateIndexes is true
	// - Start background batch processor
	// - Print initialization info
	return nil
}

// Write adds events to the batch for database insertion
func (m *MongoSink) Write(ctx context.Context, events []sinks.Event) error {
	// TODO: Implement
	// - Lock batch mutex
	// - Convert events to MongoDB documents
	// - Add events and logs to respective batches
	// - Check if batch is full and flush if needed
	// - Unlock mutex
	return nil
}

// Close cleanly shuts down the MongoDB sink
func (m *MongoSink) Close() error {
	// TODO: Implement
	// - Signal shutdown to background processor
	// - Wait for background goroutine to finish
	// - Flush any remaining events in batches
	// - Close MongoDB client connection
	// - Print final statistics
	return nil
}

// createIndexes creates database indexes for performance
func (m *MongoSink) createIndexes() error {
	// TODO: Implement
	// - Create indexes for events collection:
	//   - blockNumber, txHash, timestamp, unique(blockNumber+txHash)
	// - Create indexes for logs collection:
	//   - blockNumber, txHash, eventType, contractAddress, eventId, timestamp
	// - Use context with timeout
	return nil
}

// batchProcessor runs in background to flush batches periodically
func (m *MongoSink) batchProcessor() {
	// TODO: Implement
	// - Run in goroutine (called with go keyword)
	// - Use ticker for periodic flushing
	// - Check if flush interval has passed
	// - Flush batches if time limit reached
	// - Handle shutdown signal
	// - Call wg.Done() when finished
}

// flushBatch inserts the current batches into MongoDB
func (m *MongoSink) flushBatch() error {
	// TODO: Implement
	// - Start MongoDB session and transaction
	// - Insert events batch with InsertMany (handle duplicates)
	// - Update log documents with event IDs
	// - Insert logs batch with InsertMany
	// - Commit transaction
	// - Update metrics (totalEvents, totalLogs, totalBatches)
	// - Clear both batches
	// - Update lastFlush timestamp
	return nil
}

// updateLogEventIDs updates the event IDs in log documents
func (m *MongoSink) updateLogEventIDs(insertedIDs []interface{}) {
	// TODO: Implement
	// - Match inserted event IDs to corresponding log documents
	// - Update EventID field in log documents
}

// eventToDocument converts a sink event to MongoDB document
func (m *MongoSink) eventToDocument(event sinks.Event) EventDocument {
	// TODO: Implement
	// - Extract fields from sinks.Event
	// - Create EventDocument with proper BSON tags
	// - Set timestamp and createdAt
	return EventDocument{}
}

// logToDocument converts an event log to MongoDB document
func (m *MongoSink) logToDocument(event sinks.Event, log *types.Log) LogDocument {
	// TODO: Implement
	// - Extract topics as string array
	// - Decode event type using erc20 package
	// - Decode event-specific data (Transfer, Approval)
	// - Create LogDocument with all fields
	// - Set createdAt timestamp
	return LogDocument{}
}

// isDuplicateKeyError checks if an error is a duplicate key error
func isDuplicateKeyError(err error) bool {
	// TODO: Implement
	// - Check if error is mongo.WriteException
	// - Look for error code 11000 (duplicate key)
	return false
}

// GetEventsByBlock retrieves events for a specific block number
func (m *MongoSink) GetEventsByBlock(blockNumber uint64) ([]EventDocument, error) {
	// TODO: Implement
	// - Create filter for blockNumber
	// - Query events collection
	// - Sort by createdAt
	// - Return results
	return nil, nil
}

// GetLogsByTxHash retrieves logs for a specific transaction
func (m *MongoSink) GetLogsByTxHash(txHash string) ([]LogDocument, error) {
	// TODO: Implement
	// - Create filter for txHash
	// - Query logs collection
	// - Sort by logIndex
	// - Return results
	return nil, nil
}

// GetStatistics returns sink statistics
func (m *MongoSink) GetStatistics() map[string]interface{} {
	// TODO: Implement
	// - Return map with totalEvents, totalLogs, totalBatches, pendingEvents, pendingLogs
	return nil
}

// GetEventsByEventType retrieves events by event type with pagination
func (m *MongoSink) GetEventsByEventType(eventType string, limit int64, skip int64) ([]LogDocument, error) {
	// TODO: Implement
	// - Create filter for eventType
	// - Apply limit and skip for pagination
	// - Sort by createdAt descending
	// - Return results
	return nil, nil
}