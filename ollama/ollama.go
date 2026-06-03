package ollama

import (
	"context"

	"github.com/zeiss/prompts"
)

// DefaultURL is the default endpoint for the Ollama API.
const DefaultURL = "http://localhost:11434/v1/"

// DefaultModel is the default model for the Ollama API.
const DefaultModel = "qwen2.5:1.5b"

// Ollama is a struct that implements the Prompter interface for the Ollama API.
type Ollama struct {
	client *prompts.Client
	prompts.Unimplemented
}

var _ prompts.Promptable = (*Ollama)(nil)

// New creates a new Ollama with the given client.
func New(client *prompts.Client) *Ollama {
	return &Ollama{client: client}
}

// Do sends a prompt completion request.
func (o *Ollama) Do(ctx context.Context, prompt *prompts.Prompt) (*prompts.Completion, error) {
	res := &prompts.Completion{}

	_, err := o.client.New().Post("responses").BodyJSON(prompt).ReceiveSuccess(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// AsStream returns a stream of completion events.
func (o *Ollama) AsStream(ctx context.Context, prompt *prompts.Prompt, stream prompts.CompletionEventStream) error {
	dec := prompts.NewCompletionEventDecoder()

	_, err := o.client.New().Post("responses").BodyJSON(prompt).ResponseDecoder(dec).ReceiveSuccess(ctx, stream)
	if err != nil {
		return err
	}

	return nil
}
