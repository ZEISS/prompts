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
	client.APIKey(os.Getenv("OLLAMA_KEY"))
	provider := ollama.New(client)

	msgs := []prompts.Input{
		{
			Role: prompts.RoleSystem,
			Content: []prompts.MessageContent{
				{
					Content: prompts.MessageContentText{
						Text: "You are a helpful assistant. You answer questions short and to the point. You are concise and to the point.",
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

	prompt := prompts.NewPrompt(
		prompts.WithModel("gemma4:cloud"),
		prompts.WithInput(msgs...),
		prompts.WithStream(),
	)

	stream := provider.CompleteChunked(context.Background(), prompt)
	for event, err := range stream {
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", event.Raw.Data)
	}
}
