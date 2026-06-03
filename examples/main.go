package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zeiss/prompts"
	"github.com/zeiss/prompts/ollama"
)

// This example demonstrates how to create a completion request with a message
// It then sends the request to the API and prints the last completion content.
func main() {
	client := prompts.NewClient(os.Getenv("OLLAMA_URL"))
	provider := ollama.New(client)

	msgs := []prompts.Input{
		{
			Role: prompts.RoleSystem,
			Content: []prompts.MessageContent{
				{
					Content: prompts.MessageContentText{
						Text: "You are a helpful assistant. You answer questions to the best of your ability.",
					},
				},
			},
		},
		{
			Role: prompts.RoleUser,
			Content: []prompts.MessageContent{
				{
					Content: prompts.MessageContentText{
						Text: "What is the definition of Pi?",
					},
				},
			},
		},
	}

	stream := prompts.NewCompletionEventStream()

	go func() {
		for event := range stream {
			fmt.Println(event.Type)
		}
	}()

	prompt := prompts.NewPrompt(
		prompts.WithModel(ollama.DefaultModel),
		prompts.WithInput(msgs...),
		prompts.WithStream(),
	)

	err := provider.AsStream(context.Background(), prompt, stream)
	if err != nil {
		panic(err)
	}
}
