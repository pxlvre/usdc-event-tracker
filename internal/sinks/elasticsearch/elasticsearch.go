package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/ethereum/go-ethereum/common"
	"usdc-event-tracker/internal/logging"
	"usdc-event-tracker/internal/sinks"
)

// Config holds Elasticsearch configuration
type Config struct {
	URLs               []string
	Username           string
	Password           string
	IndexPrefix        string
	BatchSize          int
	FlushInterval      time.Duration
	UseTimestampSuffix bool // Add daily index suffix like "-2024.01.15"
}

// Sink implements the sinks.Sink interface for Elasticsearch
type Sink struct {
	config Config
	client *elasticsearch.Client
	logger *logging.Logger
}

// USDCEventDocument represents a USDC event document for Elasticsearch
type USDCEventDocument struct {
	Timestamp     string                 `json:"@timestamp"`
	BlockNumber   uint64                 `json:"block_number"`
	TxHash        string                 `json:"tx_hash"`
	TxIndex       uint                   `json:"tx_index"`
	Status        uint64                 `json:"status"`
	GasUsed       uint64                 `json:"gas_used"`
	FromAddress   string                 `json:"from_address"`
	ToAddress     string                 `json:"to_address"`
	ContractAddr  string                 `json:"contract_address"`
	Network       string                 `json:"network"`
	Events        []USDCLogEvent         `json:"events"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// USDCLogEvent represents a decoded USDC log event
type USDCLogEvent struct {
	Type         string   `json:"type"`          // Transfer, Approval, etc.
	Address      string   `json:"address"`       // Contract address
	Topics       []string `json:"topics"`        // Raw topics
	Data         string   `json:"data"`          // Raw data
	BlockNumber  uint64   `json:"block_number"`  // Block number
	TxHash       string   `json:"tx_hash"`       // Transaction hash
	TxIndex      uint     `json:"tx_index"`      // Transaction index
	LogIndex     uint     `json:"log_index"`     // Log index
	FromAddr     string   `json:"from_addr,omitempty"`     // Decoded from address
	ToAddr       string   `json:"to_addr,omitempty"`       // Decoded to address
	Value        string   `json:"value,omitempty"`         // Decoded value
	Owner        string   `json:"owner,omitempty"`         // For Approval events
	Spender      string   `json:"spender,omitempty"`       // For Approval events
}

// NewConfig creates a new Elasticsearch configuration from environment variables
func NewConfig() Config {
	config := Config{
		URLs:               []string{"http://localhost:9200"},
		IndexPrefix:        "usdc-events",
		BatchSize:          100,
		FlushInterval:      5 * time.Second,
		UseTimestampSuffix: true,
	}

	// Parse URLs from environment
	if urls := os.Getenv("ELASTICSEARCH_URLS"); urls != "" {
		config.URLs = strings.Split(urls, ",")
		for i, url := range config.URLs {
			config.URLs[i] = strings.TrimSpace(url)
		}
	}

	// Authentication
	config.Username = os.Getenv("ELASTICSEARCH_USERNAME")
	config.Password = os.Getenv("ELASTICSEARCH_PASSWORD")

	// Index configuration
	if prefix := os.Getenv("ELASTICSEARCH_INDEX_PREFIX"); prefix != "" {
		config.IndexPrefix = prefix
	}

	// Batch configuration
	if batchSize := os.Getenv("ELASTICSEARCH_BATCH_SIZE"); batchSize != "" {
		if size, err := strconv.Atoi(batchSize); err == nil && size > 0 {
			config.BatchSize = size
		}
	}

	// Timestamp suffix configuration
	if suffix := os.Getenv("ELASTICSEARCH_USE_TIMESTAMP_SUFFIX"); suffix != "" {
		config.UseTimestampSuffix = strings.ToLower(suffix) == "true"
	}

	return config
}

// New creates a new Elasticsearch sink
func New(config Config) *Sink {
	return &Sink{
		config: config,
		logger: logging.GetLogger("elasticsearch-sink"),
	}
}

// Name returns the sink name
func (s *Sink) Name() string {
	return "elasticsearch"
}

// Initialize sets up the Elasticsearch client and creates index templates
func (s *Sink) Initialize() error {
	// Create Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: s.config.URLs,
		Username:  s.config.Username,
		Password:  s.config.Password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		s.logger.Error("Failed to create Elasticsearch client", err)
		return fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	s.client = client

	// Test connection
	res, err := s.client.Info()
	if err != nil {
		s.logger.Error("Failed to connect to Elasticsearch", err)
		return fmt.Errorf("failed to connect to Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		s.logger.Error("Elasticsearch connection error", nil, map[string]interface{}{
			"status": res.Status(),
		})
		return fmt.Errorf("elasticsearch connection error: %s", res.Status())
	}

	s.logger.Info("Connected to Elasticsearch", map[string]interface{}{
		"urls":         s.config.URLs,
		"index_prefix": s.config.IndexPrefix,
		"batch_size":   s.config.BatchSize,
	})

	// Create index template for USDC events
	if err := s.createIndexTemplate(); err != nil {
		return fmt.Errorf("failed to create index template: %w", err)
	}

	return nil
}

// Write sends events to Elasticsearch
func (s *Sink) Write(ctx context.Context, events []sinks.Event) error {
	if len(events) == 0 {
		return nil
	}

	start := time.Now()
	
	// Convert events to Elasticsearch documents
	docs := s.convertEventsToDocuments(events)
	
	// Bulk index documents
	if err := s.bulkIndex(ctx, docs); err != nil {
		s.logger.Error("Failed to bulk index documents", err, map[string]interface{}{
			"event_count": len(events),
			"doc_count":   len(docs),
		})
		return fmt.Errorf("failed to bulk index documents: %w", err)
	}

	s.logger.Info("Successfully indexed events to Elasticsearch", map[string]interface{}{
		"event_count": len(events),
		"doc_count":   len(docs),
		"duration_ms": time.Since(start).Milliseconds(),
	})

	return nil
}

// Close cleans up the Elasticsearch sink
func (s *Sink) Close() error {
	s.logger.Info("Closing Elasticsearch sink")
	return nil
}

// convertEventsToDocuments converts sink events to Elasticsearch documents
func (s *Sink) convertEventsToDocuments(events []sinks.Event) []USDCEventDocument {
	docs := make([]USDCEventDocument, 0, len(events))
	
	for _, event := range events {
		// Convert logs to events
		logEvents := make([]USDCLogEvent, 0, len(event.Logs))
		for _, log := range event.Logs {
			logEvent := USDCLogEvent{
				Type:        s.decodeEventType(log.Topics),
				Address:     log.Address.Hex(),
				Topics:      s.topicsToStrings(log.Topics),
				Data:        common.Bytes2Hex(log.Data),
				BlockNumber: log.BlockNumber,
				TxHash:      log.TxHash.Hex(),
				TxIndex:     log.TxIndex,
				LogIndex:    log.Index,
			}
			
			// Decode specific event data
			s.decodeEventData(&logEvent, log.Topics, log.Data)
			
			logEvents = append(logEvents, logEvent)
		}
		
		doc := USDCEventDocument{
			Timestamp:    time.Now().UTC().Format(time.RFC3339Nano),
			BlockNumber:  event.BlockNumber,
			TxHash:       event.Receipt.TxHash.Hex(),
			TxIndex:      event.Receipt.TransactionIndex,
			Status:       event.Receipt.Status,
			GasUsed:      event.Receipt.GasUsed,
			FromAddress:  "", // Will be filled if available
			ToAddress:    "", // Will be filled if available
			ContractAddr: "", // Will be filled if available
			Network:      s.getNetworkFromConfig(),
			Events:       logEvents,
			Metadata: map[string]interface{}{
				"cumulative_gas_used": event.Receipt.CumulativeGasUsed,
				"effective_gas_price": event.Receipt.EffectiveGasPrice.String(),
				"logs_count":          len(event.Receipt.Logs),
				"usdc_logs_count":     len(event.Logs),
			},
		}
		
		docs = append(docs, doc)
	}
	
	return docs
}

// bulkIndex performs bulk indexing of documents
func (s *Sink) bulkIndex(ctx context.Context, docs []USDCEventDocument) error {
	if len(docs) == 0 {
		return nil
	}

	var buf bytes.Buffer
	indexName := s.getIndexName()

	for _, doc := range docs {
		// Bulk API format: { "index": { "_index": "indexname" } }
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
			},
		}
		
		metaBytes, _ := json.Marshal(meta)
		buf.Write(metaBytes)
		buf.WriteByte('\n')
		
		// Document data
		docBytes, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}
		buf.Write(docBytes)
		buf.WriteByte('\n')
	}

	// Perform bulk request
	req := esapi.BulkRequest{
		Body:    strings.NewReader(buf.String()),
		Refresh: "false",
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk request error: %s", res.Status())
	}

	return nil
}

// createIndexTemplate creates an index template for USDC events
func (s *Sink) createIndexTemplate() error {
	template := map[string]interface{}{
		"index_patterns": []string{s.config.IndexPrefix + "-*"},
		"template": map[string]interface{}{
			"settings": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 0,
				"refresh_interval":   "5s",
			},
			"mappings": map[string]interface{}{
				"properties": map[string]interface{}{
					"@timestamp":       map[string]interface{}{"type": "date"},
					"block_number":     map[string]interface{}{"type": "long"},
					"tx_hash":          map[string]interface{}{"type": "keyword"},
					"tx_index":         map[string]interface{}{"type": "integer"},
					"status":           map[string]interface{}{"type": "integer"},
					"gas_used":         map[string]interface{}{"type": "long"},
					"from_address":     map[string]interface{}{"type": "keyword"},
					"to_address":       map[string]interface{}{"type": "keyword"},
					"contract_address": map[string]interface{}{"type": "keyword"},
					"network":          map[string]interface{}{"type": "keyword"},
					"events": map[string]interface{}{
						"type": "nested",
						"properties": map[string]interface{}{
							"type":         map[string]interface{}{"type": "keyword"},
							"address":      map[string]interface{}{"type": "keyword"},
							"topics":       map[string]interface{}{"type": "keyword"},
							"data":         map[string]interface{}{"type": "text", "index": false},
							"block_number": map[string]interface{}{"type": "long"},
							"tx_hash":      map[string]interface{}{"type": "keyword"},
							"tx_index":     map[string]interface{}{"type": "integer"},
							"log_index":    map[string]interface{}{"type": "integer"},
							"from_addr":    map[string]interface{}{"type": "keyword"},
							"to_addr":      map[string]interface{}{"type": "keyword"},
							"value":        map[string]interface{}{"type": "keyword"},
							"owner":        map[string]interface{}{"type": "keyword"},
							"spender":      map[string]interface{}{"type": "keyword"},
						},
					},
				},
			},
		},
	}

	templateBytes, err := json.Marshal(template)
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	req := esapi.IndicesPutIndexTemplateRequest{
		Name: s.config.IndexPrefix + "-template",
		Body: bytes.NewReader(templateBytes),
	}

	res, err := req.Do(context.Background(), s.client)
	if err != nil {
		return fmt.Errorf("failed to create index template: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index template creation error: %s", res.Status())
	}

	s.logger.Info("Created Elasticsearch index template", map[string]interface{}{
		"template_name": s.config.IndexPrefix + "-template",
		"index_pattern": s.config.IndexPrefix + "-*",
	})

	return nil
}

// getIndexName returns the index name for the current time
func (s *Sink) getIndexName() string {
	if s.config.UseTimestampSuffix {
		suffix := time.Now().UTC().Format("2006.01.02")
		return fmt.Sprintf("%s-%s", s.config.IndexPrefix, suffix)
	}
	return s.config.IndexPrefix
}

// Helper methods for event decoding
func (s *Sink) decodeEventType(topics []common.Hash) string {
	if len(topics) == 0 {
		return "unknown"
	}
	
	// Transfer event signature: 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
	transferSig := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	// Approval event signature: 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925
	approvalSig := common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
	
	switch topics[0] {
	case transferSig:
		return "Transfer"
	case approvalSig:
		return "Approval"
	default:
		return "unknown"
	}
}

func (s *Sink) topicsToStrings(topics []common.Hash) []string {
	result := make([]string, len(topics))
	for i, topic := range topics {
		result[i] = topic.Hex()
	}
	return result
}

func (s *Sink) decodeEventData(event *USDCLogEvent, topics []common.Hash, data []byte) {
	// Basic decoding for Transfer and Approval events
	// In a production system, you'd use proper ABI decoding
	
	if len(topics) >= 3 {
		switch event.Type {
		case "Transfer":
			if len(topics) >= 3 {
				event.FromAddr = common.BytesToAddress(topics[1].Bytes()).Hex()
				event.ToAddr = common.BytesToAddress(topics[2].Bytes()).Hex()
			}
			if len(data) >= 32 {
				// Value is in the data field for Transfer events
				event.Value = common.BytesToHash(data[:32]).Hex()
			}
		case "Approval":
			if len(topics) >= 3 {
				event.Owner = common.BytesToAddress(topics[1].Bytes()).Hex()
				event.Spender = common.BytesToAddress(topics[2].Bytes()).Hex()
			}
			if len(data) >= 32 {
				// Value is in the data field for Approval events
				event.Value = common.BytesToHash(data[:32]).Hex()
			}
		}
	}
}

func (s *Sink) getNetworkFromConfig() string {
	// This would ideally come from the config
	// For now, we'll use an environment variable or default
	if network := os.Getenv("NETWORK"); network != "" {
		return network
	}
	return "unknown"
}