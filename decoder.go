package prompts

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"

	"github.com/katallaxie/pkg/conv"
	"github.com/katallaxie/pkg/slices"
)

const (
	// EventResponseStreamed is an event type that indicates that the response has been created.
	EventResponseCreated = "response.created"
	// EventResponseCompleted is an event type that indicates that the response has been completed.
	EventResponseCompleted = "respone.completed"
	// EventResponseFailed is an event type that indicates that the response has failed.
	EventResponseFailed = "response.failed"
)

// CompletionChunk is an event type that indicates that the response has been streamed.
type CompletionChunk struct {
	// Type is the type of the event.
	Type string `json:"type"`
	// SequenceNumber is the sequence number of the event.
	SequenceNumber int `json:"sequence_number"`
	// Data is the data associated with the event.
	Response *Completion `json:"response"`
	// Raw is the raw data associated with the event. This is useful for debugging purposes.
	Raw SSEvent `json:"-"`
}

// CompletionDecoder decodes http responses into [prompts.Completion] values.
type CompletionDecoder interface {
	// Decode decodes the response into the value pointed to by v.
	Decode(resp *http.Response) (*Completion, error)
}

// CompletionChunksDecoder decodes http responses into [prompts.Completion] values using a streaming approach.
type CompletionChunksDecoder interface {
	// Decode decodes the response into the value pointed to by v.
	Decode(resp *http.Response) iter.Seq2[*CompletionChunk, error]
}

const maxBufferSize = 512 * 1 * 1000

// ResponseDecoder decodes http responses into struct values.
type ResponseDecoder interface {
	// Decode decodes the response into the value pointed to by v.
	Decode(resp *http.Response, v any) error
}

// responsesCompatibilityDecoder decodes http responses into struct values.
type responsesCompatibilityDecoder struct{}

// NewResponsesCompatibilityDecoder returns a ResponseDecoder that decodes JSON responses into struct values.
func NewResponsesCompatibilityDecoder() responsesCompatibilityDecoder {
	return responsesCompatibilityDecoder{}
}

// Decode decodes the response to a compatible struct value.
func (d responsesCompatibilityDecoder) Decode(resp *http.Response) (*Completion, error) {
	v := &Completion{}
	err := json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// SSEvent is an event emitted by the server.
type SSEvent struct {
	// ID is the unique identifier for this event.
	ID []byte
	// Event is the name of the event.
	Event []byte
	// Data is the content of this message.
	Data []byte
	// Retry is the reconnection time (milliseconds).
	Retry []byte
}

var (
	idField    = []byte(`id`)
	eventField = []byte(`event`)
	dataField  = []byte(`data`)
	retryField = []byte(`retry`)
)

// ErrInvalidSequence is returned when the event sequence is invalid.
var ErrInvalidSequence = fmt.Errorf("invalid event sequence")

type responsesCompatibilityChunksDecoder struct{}

// NewResponsesCompatibilityChunksDecoder returns a ResponseDecoder that decodes JSON responses into struct values.
func NewResponsesCompatibilityChunksDecoder() responsesCompatibilityChunksDecoder {
	return responsesCompatibilityChunksDecoder{}
}

// Decode is decoding the responses as a stream of completion events.
func (d responsesCompatibilityChunksDecoder) Decode(resp *http.Response) iter.Seq2[*CompletionChunk, error] {
	return func(yield func(*CompletionChunk, error) bool) {
		scanner := bufio.NewScanner(resp.Body)
		buf := make([]byte, 0, maxBufferSize)
		scanner.Buffer(buf, maxBufferSize)

		for scanner.Scan() {
			sseEvent := SSEvent{}
			payload := scanner.Bytes()

			if len(payload) == 0 {
				continue
			}

			// parse line
			del := bytes.IndexByte(payload, ':')
			if del == 0 || del < 0 {
				continue // skip comment
			}

			field, content := payload[:del], bytes.TrimSpace(payload[del+1:])

			switch {
			case bytes.EqualFold(field, idField):
				sseEvent.ID = content
			case bytes.EqualFold(field, eventField):
				sseEvent.Event = content
			case bytes.EqualFold(field, dataField):
				if sseEvent.Data == nil {
					sseEvent.Data = make([]byte, 0)
				} else {
					sseEvent.Data = append(sseEvent.Data, '\n')
				}
				sseEvent.Data = append(sseEvent.Data, content...)
			case bytes.EqualFold(field, retryField):
				sseEvent.Retry = content
			}

			chunk := &CompletionChunk{
				Type:     conv.String(sseEvent.Event),
				Response: nil,
				Raw:      sseEvent,
			}

			if !slices.GreaterThen(0, sseEvent.Data...) {
				continue
			}

			c := &Completion{}
			err := json.Unmarshal(sseEvent.Data, &c)

			chunk.Response = c

			if !yield(chunk, err) {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			yield(nil, err)
		}
	}
}
