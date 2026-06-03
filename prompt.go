package prompts

import (
	"encoding/base64"
	"encoding/json"
)

type ToolChoice string

const (
	// ToolChoiceAuto is the auto tool choice.
	ToolChoiceAuto ToolChoice = "auto"
	// ToolChoiceAll is the all tool choice.
	ToolChoiceNone ToolChoice = "none"
	// ToolChoiceRequired is the required tool choice.
	ToolChoiceRequired ToolChoice = "required"
)

type isCompletionTool interface {
	isCompletionTool()
}

// Tool is an interface that represents a tool for the chat completion request.
type Tool struct {
	Tool isCompletionTool
}

func (c Tool) isCompletionTool() {}

// MarshalJSON marshals the response tool into JSON.
func (c Tool) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Tool)
}

// FunctionTool represents a function tool for the chat completion request.
type FunctionTool struct {
	// Function is the function for the chat completion request.
	Function FunctionDefinition `json:"function,omitzero"`
}

// MarshalJSON marshals the response function tool into JSON.
func (c FunctionTool) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name,omitempty"`
		// Description is the description of the function.
		Description string `json:"description,omitempty"`
		// Parameters is the parameters for the function.
		Parameters FunctionParameters `json:"parameters,omitzero"`
		// Strict is a flag to indicate whether to strictly enforce the parameters.
		Strict bool `json:"strict,omitempty"`
	}{
		Type:        "function",
		Name:        c.Function.Name,
		Description: c.Function.Description,
		Parameters:  c.Function.Parameters,
		Strict:      c.Function.Strict,
	})
}

// FunctionDefinition represents the function definition for the chat completion request.
type FunctionDefinition struct {
	// Name is the name of the function.
	Name string `json:"name"`
	// Description is the description of the function.
	Description string `json:"description,omitempty"`
	// Parameters is the parameters for the function.
	Parameters FunctionParameters `json:"parameters,omitzero"`
	// Strict is a flag to indicate whether to strictly enforce the parameters.
	Strict bool `json:"strict,omitempty"`
}

func (c FunctionTool) isCompletionTool() {}

// FunctionProperties represents the properties for the function tool.
type FunctionProperties map[string]json.RawMessage

// FunctionParameters represents the parameters for the function tool.
type FunctionParameters struct {
	// Properties is the properties for the function tool.
	Properties FunctionProperties `json:"properties,omitempty"`
	// Required is the required parameters for the function tool.
	Required []string `json:"required,omitempty"`
}

// MarshalJSON marshals the response function parameters into JSON.
func (c FunctionParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string                     `json:"type"`
		Properties map[string]json.RawMessage `json:"properties,omitempty"`
		Required   []string                   `json:"required,omitempty"`
	}{
		Type:       "object",
		Properties: c.Properties,
		Required:   c.Required,
	})
}

// CustomTool represents a custom tool for the chat completion request.
type CustomTool struct {
	// Custom is the custom tool for the chat completion request.
	Custom CustomDefinition `json:"custom,omitzero"`
}

// MarshalJSON marshals the response custom tool into JSON.
func (c CustomTool) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string           `json:"type"`
		Custom CustomDefinition `json:"custom,omitzero"`
	}{
		Type:   "custom",
		Custom: c.Custom,
	})
}

// CustomDefinition represents the custom definition for the chat completion request.
type CustomDefinition struct {
	// Name is the name of the custom tool.
	Name string `json:"name"`
	// Description is the description of the custom tool.
	Description string `json:"description,omitempty"`
}

func (c CustomTool) isCompletionTool() {}

// MessageContent is the content of a response message.
type MessageContent struct {
	Content isMessageContent
}

// MarshalJSON marshals the response message content into JSON.
func (c MessageContent) MarshalJSON() ([]byte, error) {
	if text, ok := c.GetText(); ok {
		return json.Marshal(text)
	}

	return json.Marshal(nil) // Return null if the content is not text
}

type isMessageContent interface {
	isMessageContent()
}

// NewMessageContent creates a new response message content.
func NewMessageContent() MessageContent {
	return MessageContent{}
}

// Reset resets the response message content.
func (c *MessageContent) Reset() {
	*c = MessageContent{}
}

// GetText returns the text content of the response message content.
func (c MessageContent) GetText() (MessageContentText, bool) {
	if text, ok := c.Content.(MessageContentText); ok {
		return text, true
	}

	return MessageContentText{}, false
}

// GetImage returns the image content of the response message content.
func (c MessageContent) GetImage() (MessageContentImage, bool) {
	if image, ok := c.Content.(MessageContentImage); ok {
		return image, true
	}

	return MessageContentImage{}, false
}

// GetFile returns the file content of the response message content.
func (c MessageContent) GetFile() (MessageContentFile, bool) {
	if file, ok := c.Content.(MessageContentFile); ok {
		return file, true
	}

	return MessageContentFile{}, false
}

// MessageContentText is the text content of a response message.
type MessageContentText struct {
	// Text is the text of the content.
	Text string `json:"text"`
}

// MarshalJSON marshals the response message content text into JSON.
func (c MessageContentText) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}{
		Type: "input_text",
		Text: c.Text,
	})
}

func (c MessageContentText) isMessageContent() {}

// MessageContentImage is the image content of a response message.
type MessageContentImage struct {
	Image Image `json:"image"`
}

func (c MessageContentImage) isMessageContent() {}

// MessageContentFile is the file content of a response message.
type MessageContentFile struct {
	File File `json:"file"`
}

func (c MessageContentFile) isMessageContent() {}

// File is the file for the response message content.
type File struct {
	// Name is the name of the file.
	Name string `json:"name"`
	// URL is the URL of the file.
	URL string `json:"url"`
}

// Image is a type that represents an image.
type Image struct {
	// URL is the URL of the image.
	URL string `json:"url,omitempty"`
	// Base64 is the base64 encoding of the image.
	Base64 string `json:"base64,omitempty"`
	// Name is the name of the image.
	Name string `json:"name,omitempty"`
}

// Encode encodes the image into a string.
func (i Image) Encode(data []byte) string {
	i.Base64 = base64.StdEncoding.EncodeToString(data)
	return i.Base64
}

// NewImage creates a new image from the given data.
func NewImage(data []byte) Image {
	var img Image
	img.Encode(data)

	return img
}

// Input is the message for chat completion.
type Input struct {
	// Role is the role of the message sender.
	Role string `json:"role"`
	// Content is the content of the message.
	Content []MessageContent `json:"content"`
	// Name is the name of the message sender (optional).
	Name string `json:"name,omitempty"`
}

// Prompt is a chat completion request.
type Prompt struct {
	// Model is the model for the chat completion request.
	Model string `json:"model"`
	// Input is the list of messages for the chat completion request.
	Input []Input `json:"input"`
	// Instructions is the instructions for the chat completion request.
	Instructions string `json:"instructions,omitempty"`
	// Tools is the list of tools to use for the chat completion request.
	Tools []Tool `json:"tools,omitempty"`
	// ToolChoice is the tool choice for the chat completion request.
	ToolChoice ToolChoice `json:"tool_choice,omitempty"`
	// MaxTokens is the maximum number of tokens for the chat completion request.
	MaxTokens *int `json:"max_tokens,omitzero"`
	// Temperature is the sampling temperature
	Temperature *float32 `json:"temperature,omitzero"`
	// Stream is a flag to enable streaming
	Stream bool `json:"stream,omitempty"`
	// TopP is the nucleus sampling parameter
	TopP *float64 `json:"top_p,omitzero"`
	// TopK is the number of top tokens to sample from
	TopK *int `json:"top_k,omitzero"`
}

// PromptOpt is a function type for configuring the ResponseRequest.
type PromptOpt func(*Prompt)

// NewPrompt creates a new chat completion request with the given options.
func NewPrompt(opts ...PromptOpt) *Prompt {
	req := new(Prompt)

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// WithInput sets the messages for the chat completion request.
func WithInput(msgs ...Input) PromptOpt {
	return func(req *Prompt) {
		req.Input = msgs
	}
}

// WithInstructions sets the instructions for the chat completion request.
func WithInstructions(instructions string) PromptOpt {
	return func(req *Prompt) {
		req.Instructions = instructions
	}
}

// WithTools sets the tools for the chat completion request.
func WithTools(tools ...Tool) PromptOpt {
	return func(req *Prompt) {
		req.Tools = tools
	}
}

// WithModel sets the model for the chat completion request.
func WithModel(model string) PromptOpt {
	return func(req *Prompt) {
		req.Model = model
	}
}

// WithStream sets whether or not to stream the chat completion response.
func WithStream() PromptOpt {
	return func(req *Prompt) {
		req.Stream = true
	}
}
