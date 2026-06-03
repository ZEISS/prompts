package prompts

import (
	"encoding/json"

	"github.com/katallaxie/pkg/utilx"
)

const (
	// FinishReasonStop indicates that the chat completion was finished because the model stopped generating content.
	FinishReasonStop = "stop"
	// FinishReasonLength indicates that the chat completion was finished because the maximum length was reached.
	FinishReasonLength = "length"
	// FinishReasonContentFilter indicates that the chat completion was finished because the content filter was triggered.
	FinishReasonContentFilter = "content_filter"
	// FinishReasonUnknown indicates that the chat completion was finished for an unknown reason.
	FinishReasonUnknown = ""
)

// ChatCompletionMessageToolCall represents a tool call in a chat completion message.
type ChatCompletionMessageToolCall struct {
	// ToolCall is the name of the tool being called.
	ToolCall isToolCall
}

type isToolCall interface {
	isToolCall()
}

// ChatCompletionMessageFunctionToolCall represents a function tool call in a chat completion message.
type ChatCompletionMessageFunctionToolCall struct {
	// ID is the unique identifier for the function tool call.
	ID string `json:"id,omitempty"`
	// Type is the type of the custom tool being called.
	Type string `json:"type,omitempty"`
	// Function is the function being called.
	Function ChatCompletionMessageFunction `json:"function,omitzero"`
}

func (ChatCompletionMessageFunctionToolCall) isToolCall() {}

// UnmarshalJSON implements the json.Unmarshaler interface for ChatCompletionMessageToolCall.
func (c *ChatCompletionMessageToolCall) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID       string                         `json:"id,omitempty"`
		Type     string                         `json:"type,omitempty"`
		Function *ChatCompletionMessageFunction `json:"function,omitempty"`
		Custom   *ChatCompletionMessageCustom   `json:"custom,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.ToolCall = nil

	if utilx.NotNil(aux.Function) {
		c.ToolCall = ChatCompletionMessageFunctionToolCall{
			ID:       aux.ID,
			Type:     aux.Type,
			Function: *aux.Function,
		}
	}

	if utilx.NotNil(aux.Custom) {
		c.ToolCall = ChatCompletionMessageCustomToolCall{
			ID:     aux.ID,
			Type:   aux.Type,
			Custom: *aux.Custom,
		}
	}

	return nil
}

// ChatCompletionMessageFunction represents a function in a chat completion message.
type ChatCompletionMessageFunction struct {
	// Name is the name of the function.
	Name string `json:"name,omitempty"`
	// Arguments is the arguments for the function.
	Arguments map[string]any `json:"arguments,omitempty"`
}

// ChatCompletionMessageCustomToolCall represents a custom tool call in a chat completion message.
type ChatCompletionMessageCustomToolCall struct {
	// ID is the unique identifier for the custom tool call.
	ID string `json:"id,omitempty"`
	// Type is the type of the custom tool being called.
	Type string `json:"type,omitempty"`
	// Custom is the custom tool being called.
	Custom ChatCompletionMessageCustom `json:"custom"`
}

func (ChatCompletionMessageCustomToolCall) isToolCall() {}

type ChatCompletionMessageCustom struct {
	// Name is the name of the custom tool.
	Name string `json:"name,omitempty"`
	// Input is the input for the custom tool.
	Input map[string]any `json:"input,omitempty"`
}

// ChatCompletionChoiceIndex is the index for the chat completion.
type ChatCompletionChoiceIndex struct {
	// Role is the role of the message sender.
	Role string `json:"role"`
	// Content is the content of the message.
	Content string `json:"content"`
	// Annotations is the annotations for the message.
	Annotations []ChatCompletionAnnotation `json:"annotations,omitempty"`
	// ToolCalls is the tool calls for the message.
	ToolCalls []ChatCompletionMessageToolCall `json:"tool_calls,omitempty"`
}

// ChatCompletionAnnotation is the annotation for the chat completion.
type ChatCompletionAnnotation struct {
	// Type is the type of the annotation.
	Type string `json:"type,omitempty"`
	// URLCitation is the URL citation for the chat completion.
	URLCitation ChatCompletionAnnotationUrlCitation `json:"url_citation,omitzero"`
}

// ChatCompletionAnnotationUrlCitation is the URL citation for the chat completion.
type ChatCompletionAnnotationUrlCitation struct {
	// Title is the title of the URL citation.
	Title string `json:"title,omitempty"`
	// URL is the URL of the URL citation.
	URL string `json:"url,omitempty"`
	// StartIndex is the start index of the URL citation in the content.
	StartIndex int `json:"start_index,omitempty"`
	// EndIndex is the end index of the URL citation in the content.
	EndIndex int `json:"end_index,omitempty"`
}

// CompletionUsage represents the usage of the chat completion.
type CompletionUsage struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens,omitempty"`
	// CompletionTokens is the number of tokens in the completion.
	CompletionTokens int `json:"completion_tokens,omitempty"`
	// TotalTokens is the total number of tokens used in the chat completion.
	TotalTokens int `json:"total_tokens,omitempty"`
}

const (
	// StatusInProgress indicates that the response is in progress.
	StatusInProgress = "in_progress"
	// StatusCompleted indicates that the response is completed.
	StatusCompleted = "completed"
	// StatusIncomplete indicates that the response is incomplete.
	StatusIncomplete = "incomplete"
)

// Output represents the output of the chat completion response.
type Output struct {
	// Output is the output of the completion of the response.
	Output isOutput
}

type isOutput interface {
	isOutput()
}

// GetFunctionCall returns the function call if the output is a function call.
func (r *Output) GetFunctionCall() (*OutputFunctionCall, bool) {
	if call, ok := r.Output.(*OutputFunctionCall); ok {
		return call, true
	}

	return nil, false
}

// GetMessage returns the message if the output is a message.
func (r *Output) GetMessage() (*OutputMessage, bool) {
	if msg, ok := r.Output.(*OutputMessage); ok {
		return msg, true
	}

	return nil, false
}

func (r *Output) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type string `json:"type,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.Output = nil

	switch aux.Type {
	case "function_call":
		var fnCall OutputFunctionCall
		if err := json.Unmarshal(data, &fnCall); err != nil {
			return err
		}
		r.Output = fnCall
	case "message":
		var msg OutputMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		r.Output = msg
	}

	return nil
}

// OutputFunctionCall represents a function call output in the chat completion response.
type OutputFunctionCall struct {
	// ID is the unique identifier for the function call output.
	ID string `json:"id,omitempty"`
	// CallID is the unique identifier for the function call.
	CallID string `json:"call_id,omitempty"`
	// Status is the status of the message output.
	Status string `json:"status,omitempty"`
	// Role is the role of the message sender.
	Role string `json:"role"`
	// Name is the name of the function being called.
	Name string `json:"name,omitempty"`
	// Arguments is the arguments for the function being called.
	Arguments string `json:"arguments,omitempty"`
}

func (OutputFunctionCall) isOutput() {}

// OutputMessage represents a message output in the chat completion response.
type OutputMessage struct {
	// ID is the unique identifier for the message output.
	ID string `json:"id,omitempty"`
	// Status is the status of the message output.
	Status string `json:"status,omitempty"`
	// Role is the role of the message sender.
	Role string `json:"role"`
	// OutputMessageContent is the content of the message output.
	OutputMessageContent []OutputMessageContent `json:"content,omitzero"`
}

func (OutputMessage) isOutput() {}

// OutputMessageContent represents the content of a message output in the chat completion response.
type OutputMessageContent struct {
	// Content is the content of the message output.
	Content isOutputMessageContent
}

type isOutputMessageContent interface {
	isOutputMessageContent()
}

// UnmarshalJSON implements the json.Unmarshaler interface for OutputMessageContent.
func (c *OutputMessageContent) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type string `json:"type,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.Content = nil

	if aux.Type == "output_text" {
		var text OutputMessageContentText
		if err := json.Unmarshal(data, &text); err != nil {
			return err
		}
		c.Content = text
	}

	return nil
}

func (OutputMessageContent) isOutputMessageContent() {}

// OutputMessageContentText represents a text content of a message output in the chat completion response.
type OutputMessageContentText struct {
	// Text is the text content of the message output.
	Text string `json:"text,omitempty"`
}

func (OutputMessageContentText) isOutputMessageContent() {}

// Completion represents a chat completion response.
type Completion struct {
	// ID is the unique identifier for the response
	ID string `json:"id,omitempty"`
	// Object is the type of object returned
	Object string `json:"object,omitempty"`
	// CreatedAt is the timestamp of when the response was created
	CreatedAt int64 `json:"created_at,omitempty"`
	// Status is the status of the response
	Status string `json:"status,omitempty"`
	// CompletedAt is the timestamp of when the response was completed
	CompletedAt int64 `json:"completed_at,omitempty"`
	// Instructions is the instructions for the chat completion response
	Instructions string `json:"instructions,omitempty"`
	// ParallelToolCalls indicates whether the tool calls were executed in parallel
	ParallelToolCalls bool `json:"parallel_tool_calls,omitempty"`
	// Output is the output of the chat completion response
	Output []Output `json:"output,omitempty"`
}

// SearchResult represents a search result structure for chat completion API.
type SearchResult struct {
	// Title is the title of the search result
	Title string `json:"title,omitempty"`
	// URL is the URL of the search result
	URL string `json:"url,omitempty"`
	// Snippet is the snippet of the search result
	Snippet string `json:"snippet,omitempty"`
	// Source is the source of the search result
	Source string `json:"source,omitempty"`
}
