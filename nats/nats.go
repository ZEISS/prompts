package nats

import (
	"context"
	"fmt"
	"iter"

	"github.com/nats-io/nats.go"
	"github.com/zeiss/prompts"
)

const (
	// DefaultMaxPayloadSize is the default maximum payload size in bytes for a Nats request.
	DefaultMaxPayloadSize = 3 * 1024
)

// ErrUnimplemented is returned when a method is not implemented.
var ErrUnimplemented = fmt.Errorf("prompts: unimplemented")

var _ prompts.Promptable = (*Nats)(nil)

// Nats is a struct that holds a NATS connection and options for configuring a NATS client.
type Nats struct {
	nc   *nats.Conn
	opts *NatsOpts
}

// NatsOpts is a struct that holds options for configuring a Nats client.
type NatsOpts struct {
	// MaxPayloadSize is the maximum payload size in bytes for a Nats request.
	MaxPayloadSize int `json:"max_payload_size"`
}

// NatsOpt is a function that configures a Client.
type NatsOpt func(*NatsOpts)

// DefaultNatsOpts returns the default options for configuring a Nats client.
func DefaultNatsOpts() *NatsOpts {
	return &NatsOpts{
		MaxPayloadSize: DefaultMaxPayloadSize,
	}
}

// WithMaxPayloadSizeOpt sets the maximum payload size in bytes for a Nats request.
func WithMaxPayloadSizeOpt(size int) NatsOpt {
	return func(opts *NatsOpts) {
		opts.MaxPayloadSize = size
	}
}

// New creates a new Nats with the given client.
func New(nc *nats.Conn, agent *AgentInfo, opts ...NatsOpt) *Nats {
	options := DefaultNatsOpts()

	for _, opt := range opts {
		opt(options)
	}

	nats := new(Nats)
	nats.nc = nc
	nats.opts = options

	return nats
}

// Complete sends a prompt and returns the completion.
func (n *Nats) Complete(ctx context.Context, prompt *prompts.Prompt) (*prompts.Completion, error) {
	return nil, ErrUnimplemented
}

// CompleteChunked sends a prompt and returns a stream of completion events.
func (n *Nats) CompleteChunked(ctx context.Context, prompt *prompts.Prompt) iter.Seq2[*prompts.CompletionChunk, error] {
	return nil
}
