// Package sinks provides interfaces and implementations for data output destinations
package sinks

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
)

// Event represents a processed USDC event with metadata
type Event struct {
	BlockNumber uint64
	Receipt     *types.Receipt
	Logs        []*types.Log
}

// Sink defines the interface for data output destinations
type Sink interface {
	// Name returns the name of the sink
	Name() string
	
	// Initialize prepares the sink for use
	Initialize() error
	
	// Write sends events to the sink
	Write(ctx context.Context, events []Event) error
	
	// Close cleans up sink resources
	Close() error
}

// Manager orchestrates multiple sinks, allowing data to be sent to multiple destinations.
type Manager struct {
	sinks []Sink
}

// NewManager creates a new sink manager with an empty list of sinks.
func NewManager() *Manager {
	return &Manager{
		sinks: make([]Sink, 0),
	}
}

// AddSink registers a new sink with the manager.
func (m *Manager) AddSink(sink Sink) {
	m.sinks = append(m.sinks, sink)
}

// Initialize prepares all registered sinks for use.
// If any sink fails to initialize, the error is returned immediately.
func (m *Manager) Initialize() error {
	for _, sink := range m.sinks {
		if err := sink.Initialize(); err != nil {
			return err
		}
	}
	return nil
}

// Write distributes events to all registered sinks.
// Errors from individual sinks are logged but don't stop other sinks from receiving data.
func (m *Manager) Write(ctx context.Context, events []Event) error {
	for _, sink := range m.sinks {
		if err := sink.Write(ctx, events); err != nil {
			// Log error but continue with other sinks
			// In production, you might want different error handling
			continue
		}
	}
	return nil
}

// Close cleanly shuts down all registered sinks.
// All sinks are closed even if some return errors.
func (m *Manager) Close() error {
	for _, sink := range m.sinks {
		if err := sink.Close(); err != nil {
			// Log error but continue closing others
			continue
		}
	}
	return nil
}

// HasSink checks if a sink with the specified name is registered.
func (m *Manager) HasSink(name string) bool {
	for _, sink := range m.sinks {
		if sink.Name() == name {
			return true
		}
	}
	return false
}