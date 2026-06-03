package prompts

import (
	"context"
	"encoding/json"
	"fmt"
)

// ErrNotImplemented is returned when a method is not implemented.
var ErrNotImplemented = fmt.Errorf("prompts: not implemented")

// NewPromptError returns a new instance of prompts.PromptError.
func NewPromptError() *PromptError {
	return &PromptError{}
}

// PromptError is an error that can be returned by the prompts package.
type PromptError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
	Type    string `json:"type"`
	Param   string `json:"param"`
}

// Error returns the error message associated with the Err instance.
func (e *PromptError) Error() string {
	return fmt.Sprint(e.Message)
}

// UnmarshalJSON unmarshals the JSON data into the Err instance.
func (e *PromptError) UnmarshalJSON(data []byte) error {
	err := struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Type    string `json:"type"`
			Param   string `json:"param"`
		} `json:"error"`
	}{}

	if err := json.Unmarshal(data, &err); err != nil {
		return err
	}

	e.Code = err.Error.Code
	e.Message = err.Error.Message
	e.Type = err.Error.Type
	e.Param = err.Error.Param

	return nil
}

// Promptable is an interface that enables to send a prompt.
type Promptable interface {
	// Do completes a prompt and returns the completion.
	Do(ctx context.Context, prompt *Prompt) (*Completion, error)
	// AsStream completes a prompt and returns the completion as a stream.
	AsStream(ctx context.Context, prompt *Prompt, stream CompletionEventStream) error
}

var _ Promptable = &Unimplemented{}

// Unimplemented is an unimplemented Promptable.
type Unimplemented struct{}

// Complete is an unimplemented method.
func (u *Unimplemented) Do(_ context.Context, _ *Prompt) (*Completion, error) {
	return nil, ErrNotImplemented
}

// AsStream is an unimplemented method.
func (u *Unimplemented) AsStream(_ context.Context, _ *Prompt, _ CompletionEventStream) error {
	return ErrNotImplemented
}
