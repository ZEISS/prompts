package prompts

const (
	// EventResponseStreamed is an event type that indicates that the response has been created.
	EventResponseCreated = "response.created"
	// EventResponseCompleted is an event type that indicates that the response has been completed.
	EventResponseCompleted = "respone.completed"
	// EventResponseFailed is an event type that indicates that the response has failed.
	EventResponseFailed = "response.failed"
)

// CompletionEvent is an event that is returned from an API.
type CompletionEvent struct {
	// Type is the type of the event.
	Type string `json:"type"`
	// SequenceNumber is the sequence number of the event.
	SequenceNumber int `json:"sequence_number"`
	// Data is the data associated with the event.
	Response *Completion `json:"response"`
	// Raw is the raw data associated with the event. This is useful for debugging purposes.
	Raw SSEvent `json:"-"`
}

// CompletionEventStream is a stream of completion events.
type CompletionEventStream chan CompletionEvent

// Close closes the stream.
func (s CompletionEventStream) Close() { close(s) }

// NewCompletionEventStream creates a new completion event stream.
func NewCompletionEventStream() CompletionEventStream {
	return make(CompletionEventStream)
}
