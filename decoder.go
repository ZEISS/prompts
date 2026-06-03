package prompts

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/katallaxie/pkg/conv"
	"github.com/katallaxie/pkg/slices"
)

const maxBufferSize = 512 * 1 * 1000

// ResponseDecoder decodes http responses into struct values.
type ResponseDecoder interface {
	// Decode decodes the response into the value pointed to by v.
	Decode(resp *http.Response, v any) error
}

// jsonDecoder decodes http response JSON into a JSON-tagged struct value.
type jsonDecoder struct{}

// NewJSONDecoder returns a ResponseDecoder that decodes JSON responses into struct values.
func NewJSONDecoder() ResponseDecoder {
	return jsonDecoder{}
}

// Decode decodes the Response Body into the value pointed to by v.
// Caller must provide a non-nil v and close the resp.Body.
func (d jsonDecoder) Decode(resp *http.Response, v any) error {
	return json.NewDecoder(resp.Body).Decode(v)
}

type byteStreamer struct{}

// NewByteStreamer returns a ResponseDecoder that copies the response body into an [io.Writer] instance.
func NewByteStreamer() ResponseDecoder {
	return byteStreamer{}
}

// Decode simply tries to copy response data into v assuming its an [io.Writer] instance. Assuming so little about v
// gives consumers a lot of choice about consuming response data. They can wait for all data to be dumped into some
// buffer then act on it or they can read as soon as data gets written.
func (d byteStreamer) Decode(resp *http.Response, v any) error {
	w, ok := v.(io.Writer)
	if !ok {
		return fmt.Errorf("expected type: %T; got: %T", w, v)
	}

	_, err := io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("failed copying response data to v: %w", err)
	}

	return nil
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

type completionEventDecoder struct{}

// NewCompletionEventDecoder returns a new SSE decoder.
func NewCompletionEventDecoder() ResponseDecoder {
	return &completionEventDecoder{}
}

// Decode simply tries to copy response data into v assuming its an [io.Writer] instance. Assuming so little about v
// gives consumers a lot of choice about consuming response data. They can wait for all data to be dumped into some
// buffer then act on it or they can read as soon as data gets written.
func (d completionEventDecoder) Decode(resp *http.Response, v any) error {
	s, ok := v.(CompletionEventStream)
	if !ok {
		return fmt.Errorf("expected type: %T; got: %T", s, v)
	}

	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)

	defer close(s)

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

		event := CompletionEvent{Raw: sseEvent}
		event.Type = conv.String(sseEvent.Event)

		if !slices.GreaterThen(0, sseEvent.Data...) {
			s <- event
			continue
		}

		err := json.Unmarshal(sseEvent.Data, &event.Response)
		if err != nil {
			return err
		}

		s <- event
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return nil
}
