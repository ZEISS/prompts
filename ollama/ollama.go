package ollama

import (
	"context"
	"iter"

	"github.com/zeiss/prompts"
)

// DefaultURL is the default endpoint for the Ollama API.
const DefaultURL = "http://localhost:11434/v1/"

// DefaultModel is the default model for the Ollama API.
const DefaultModel = "qwen2.5:1.5b"

var _ prompts.Promptable = (*Ollama)(nil)

var (
	CompletionDecoder       = prompts.NewResponsesCompatibilityDecoder()
	CompletionChunksDecoder = prompts.NewResponsesCompatibilityChunksDecoder()
)

// Ollama is a struct that implements the Prompter interface for the Ollama API.
type Ollama struct {
	client *prompts.Client
}

// New creates a new Ollama with the given client.
func New(client *prompts.Client) *Ollama {
	return &Ollama{client: client}
}

// Complete sends a prompt completion request.
func (o *Ollama) Complete(ctx context.Context, prompt *prompts.Prompt) (*prompts.Completion, error) {
	return o.client.New().Post("responses").Complete(ctx, prompt)
}

// CompleteChunked sends a prompt and returns a stream of completion events.
func (o *Ollama) CompleteChunked(ctx context.Context, prompt *prompts.Prompt) iter.Seq2[*prompts.CompletionChunk, error] {
	return o.client.New().Post("responses").CompleteChunked(ctx, prompt)
}
