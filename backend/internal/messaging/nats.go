package messaging

import (
	"backend/internal/config"
	"fmt"

	"github.com/nats-io/nats.go"
)

// NATSClient wraps nats.Conn
type NATSClient struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

// New creates a new NATS client
func New(cfg config.NATSConfig) (*NATSClient, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &NATSClient{
		conn: conn,
		js:   js,
	}, nil
}

// GetConn returns the underlying NATS connection
func (n *NATSClient) GetConn() *nats.Conn {
	return n.conn
}

// GetJetStream returns the JetStream context
func (n *NATSClient) GetJetStream() nats.JetStreamContext {
	return n.js
}

// Close closes the NATS connection
func (n *NATSClient) Close() {
	n.conn.Close()
}
