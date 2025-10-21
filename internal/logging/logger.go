package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

type LogEntry struct {
	Timestamp string                 `json:"@timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Component string                 `json:"component"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

type Logger struct {
	component string
	enabled   bool
}

var globalLogger *Logger

func Init(component string) {
	globalLogger = &Logger{
		component: component,
		enabled:   true,
	}
}

func GetLogger(component string) *Logger {
	return &Logger{
		component: component,
		enabled:   true,
	}
}

func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	if !l.enabled {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Message:   message,
		Component: l.component,
		Fields:    fields,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	fmt.Fprintln(os.Stdout, string(jsonBytes))
}

func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(DEBUG, message, f)
}

func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(INFO, message, f)
}

func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(WARN, message, f)
}

func (l *Logger) Error(message string, err error, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	} else {
		f = make(map[string]interface{})
	}
	
	if err != nil {
		f["error"] = err.Error()
	}
	
	l.log(ERROR, message, f)
}

// Convenience functions for common blockchain events
func (l *Logger) LogBlockProcessing(blockNumber uint64, txCount int) {
	l.Info("Processing blockchain block", map[string]interface{}{
		"block_number":      blockNumber,
		"transaction_count": txCount,
		"event_type":        "block_processing",
	})
}

func (l *Logger) LogUSDCTransaction(txHash string, blockNumber uint64, gasUsed uint64, eventType string, fromAddr, toAddr string, value string) {
	l.Info("USDC transaction detected", map[string]interface{}{
		"tx_hash":       txHash,
		"block_number":  blockNumber,
		"gas_used":      gasUsed,
		"event_type":    eventType,
		"from_address":  fromAddr,
		"to_address":    toAddr,
		"value":         value,
		"currency":      "USDC",
		"event_category": "blockchain_event",
	})
}

func (l *Logger) LogSinkOperation(sinkName string, operation string, eventCount int, duration time.Duration, success bool) {
	fields := map[string]interface{}{
		"sink_name":    sinkName,
		"operation":    operation,
		"event_count":  eventCount,
		"duration_ms":  duration.Milliseconds(),
		"success":      success,
		"event_type":   "sink_operation",
	}
	
	if success {
		l.Info("Sink operation completed", fields)
	} else {
		l.Warn("Sink operation failed", fields)
	}
}

func (l *Logger) LogConnection(network string, chainID string, usdcAddress string, webhookURL string) {
	l.Info("Ethereum connection established", map[string]interface{}{
		"network":      network,
		"chain_id":     chainID,
		"usdc_address": usdcAddress,
		"webhook_url":  webhookURL,
		"event_type":   "connection_established",
	})
}

func (l *Logger) LogError(message string, err error, fields ...map[string]interface{}) {
	l.Error(message, err, fields...)
}

// Global convenience functions
func Debug(message string, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(message, fields...)
	}
}

func Info(message string, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Info(message, fields...)
	}
}

func Warn(message string, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(message, fields...)
	}
}

func Error(message string, err error, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Error(message, err, fields...)
	}
}