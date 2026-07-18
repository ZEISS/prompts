package nats

import (
	"encoding/json"
)

// PromptMessage is a message sent to a NATS JetStream agent.
type PromptMessage struct {
	// Messages is a list of messages to be sent with the request.
	Messages []Message `json:"messages"`
	// Attachments is a list of attachments to be sent with the request.
	Attachments []Attachment `json:"attachments"`
	// Models is a list of models to be used with the request.
	Models []Model `json:"models"`
}

// Attachment is an attachment to a request message.
type Attachment struct {
	// Filename is the name of the file to be attached.
	Filename string `json:"filename"`
	// FileContent is the content of the file to be attached.
	// The content is base64 encoded.
	Content []byte `json:"content"`
}

// Model is the model of the request message.
type Model struct {
	// Name is the name of the model to be used.
	Name string `json:"name"`
	// Parameters is a map of parameters to be used with the model.
	Parameters map[string]string `json:"parameters"`
	// Priority is the priority of the model to be used.
	Priority int `json:"priority"`
}

// Message is a message sent to a NATS JetStream agent.
type Message struct {
	// Role is the role of the message sender.
	Role string `json:"role"`
	// Content is the content of the message.
	Content string `json:"content"`
}

// Tool is a tool to be used with a request message.
type Tool struct {
	// Type is the type of the tool to be used.
	Type string `json:"type"`
	// Name is the name of the tool to be used.
	Name string `json:"name"`
	// Description is the description of the tool to be used.
	Description string `json:"description"`
	// Parameters is a map of parameters to be used with the tool.
	Parameters ToolParameters `json:"parameters"`
}

// ToolParameters is a map of parameters for a tool to be used with a request message.
type ToolParameters map[string]ToolParameter

// ToolParameter is a parameter for a tool to be used with a request message.
type ToolParameter struct {
	// Name is the name of the parameter.
	Name string `json:"name"`
	// Description is the description of the parameter.
	Description string `json:"description"`
	// Type is the type of the parameter.
	Type string `json:"type"`
	// AdditionalProperties is a map of additional properties for the parameter.
	AdditionalProperties bool
}

// RawChunk is a raw chunk of data received from a NATS JetStream agent.
type RawChunk struct {
	// Type is the type of the chunk.
	Type string `json:"type"`
	// Data is the raw data of the chunk.
	Data json.RawMessage `json:"data"`
}

// QueryChunk is a chunk of data received from a NATS JetStream agent as a query response.
type QueryChunk struct {
	// Type is the type of the chunk.
	Type string `json:"type"`
	// Data is the raw data of the chunk.
	Data struct {
		// ID is the ID of the query response.
		ID string `json:"id"`
		// ReplySubject is the subject to which the reply should be sent.
		ReplySubject string `json:"reply_subject"`
		// Prompt is the prompt that was used to generate the response.
		Prompt string `json:"prompt"`
	} `json:"data"`
}
