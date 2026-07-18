package prompts

import (
	"context"
	"encoding/base64"
	"io"
	"iter"
	"maps"
	"net/http"
	"net/url"
	"time"

	goquery "github.com/google/go-querystring/query"
	"github.com/katallaxie/pkg/slices"
)

const (
	contentType      = "Content-Type"
	jsonContentType  = "application/json"
	streamAcceptType = "text/event-stream"
	formContentType  = "application/x-www-form-urlencoded"
)

// DefaultTimeout is the default timeout for the Perplexity API.
const DefaultTimeout = 30 * time.Second

// DefaultClient is the default HTTP client for the Perplexity API.
var DefaultClient = &http.Client{
	Timeout: DefaultTimeout,
}

// Doer executes http requests. It is implemented by *http.Client.  You can
// wrap *http.Client with layers of Doers to form a stack of client-side
// middleware.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

var _ Promptable = (*Client)(nil)

// Client is a struct for sending HTTP requests. It holds the http Client and
// the response decoder.
type Client struct {
	// http Client for doing requests
	httpClient Doer
	// HTTP method (GET, POST, etc.)
	method string
	// raw url string for requests
	rawURL string
	// stores key-values pairs to add to request's Headers
	header http.Header
	// url tagged query structs
	queryStructs []any
	// body provider
	bodyProvider BodyProvider
	// response decoder
	responseDecoder ResponseDecoder
	// completion decoder
	completionDecoder CompletionDecoder
	// completion chunks decoder
	completionChunksDecoder CompletionChunksDecoder
}

// NewClient returns a new Client with an http DefaultClient.
func NewClient(url ...string) *Client {
	c := &Client{
		httpClient:              DefaultClient,
		method:                  http.MethodGet,
		header:                  make(http.Header),
		queryStructs:            make([]any, 0),
		bodyProvider:            jsonBodyProvider{},
		completionDecoder:       responsesCompatibilityDecoder{},
		completionChunksDecoder: responsesCompatibilityChunksDecoder{},
	}

	if slices.GreaterThen(0, url) {
		c.rawURL = slices.First(url...)
	}

	return c
}

// New returns a copy of a Client for creating a new Client with properties
// from a parent Client. For example,
//
//	parentClient := clients.New().Client(client).Base("https://api.io/")
//	fooClient := parentClient.New().Get("foo/")
//	barClient := parentClient.New().Get("bar/")
//
// fooClient and barClient will both use the same client, but send requests to
// https://api.io/foo/ and https://api.io/bar/ respectively.
//
// Note that query and body values are copied so if pointer values are used,
// mutating the original value will mutate the value within the child Client.
func (s *Client) New() *Client {
	// copy Headers pairs into new Header map
	headerCopy := make(http.Header)
	maps.Copy(headerCopy, s.header)

	return &Client{
		httpClient:              s.httpClient,
		method:                  s.method,
		rawURL:                  s.rawURL,
		header:                  headerCopy,
		queryStructs:            append([]any{}, s.queryStructs...),
		bodyProvider:            s.bodyProvider,
		completionDecoder:       s.completionDecoder,
		completionChunksDecoder: s.completionChunksDecoder,
	}
}

// Client sets the http Client used to do requests. If a nil client is given,
// the http.DefaultClient will be used.
func (s *Client) Client(httpClient *http.Client) *Client {
	if httpClient == nil {
		return s.Doer(http.DefaultClient)
	}

	return s.Doer(httpClient)
}

// Doer sets the custom Doer implementation used to do requests.
// If a nil client is given, the http.DefaultClient will be used.
func (s *Client) Doer(doer Doer) *Client {
	if doer == nil {
		s.httpClient = http.DefaultClient
	} else {
		s.httpClient = doer
	}

	return s
}

// Method

// Head sets the Client method to HEAD and sets the given pathURL.
func (s *Client) Head(pathURL string) *Client {
	s.method = http.MethodHead
	return s.Path(pathURL)
}

// Get sets the Client method to GET and sets the given pathURL.
func (s *Client) Get(pathURL string) *Client {
	s.method = http.MethodGet
	return s.Path(pathURL)
}

// Post sets the Client method to POST and sets the given pathURL.
func (s *Client) Post(pathURL string) *Client {
	s.method = http.MethodPost
	return s.Path(pathURL)
}

// Put sets the Client method to PUT and sets the given pathURL.
func (s *Client) Put(pathURL string) *Client {
	s.method = http.MethodPut
	return s.Path(pathURL)
}

// Patch sets the Client method to PATCH and sets the given pathURL.
func (s *Client) Patch(pathURL string) *Client {
	s.method = http.MethodPatch
	return s.Path(pathURL)
}

// Delete sets the Client method to DELETE and sets the given pathURL.
func (s *Client) Delete(pathURL string) *Client {
	s.method = http.MethodDelete
	return s.Path(pathURL)
}

// Options sets the Client method to OPTIONS and sets the given pathURL.
func (s *Client) Options(pathURL string) *Client {
	s.method = http.MethodOptions
	return s.Path(pathURL)
}

// Trace sets the Client method to TRACE and sets the given pathURL.
func (s *Client) Trace(pathURL string) *Client {
	s.method = http.MethodTrace
	return s.Path(pathURL)
}

// Connect sets the Client method to CONNECT and sets the given pathURL.
func (s *Client) Connect(pathURL string) *Client {
	s.method = http.MethodConnect
	return s.Path(pathURL)
}

// Header

// Add adds the key, value pair in Headers, appending values for existing keys
// to the key's values. Header keys are canonicalized.
func (s *Client) Add(key, value string) *Client {
	s.header.Add(key, value)
	return s
}

// Set sets the key, value pair in Headers, replacing existing values
// associated with key. Header keys are canonicalized.
func (s *Client) Set(key, value string) *Client {
	s.header.Set(key, value)
	return s
}

// APIKey sets the Authorization header to use the provided API key with the Bearer scheme.
func (s *Client) APIKey(apiKey string) *Client {
	return s.Set("Authorization", "Bearer "+apiKey)
}

// SetBasicAuth sets the Authorization header to use HTTP Basic Authentication
// with the provided username and password. With HTTP Basic Authentication
// the provided username and password are not encrypted.
func (s *Client) SetBasicAuth(username, password string) *Client {
	return s.Set("Authorization", "Basic "+basicAuth(username, password))
}

// basicAuth returns the base64 encoded username:password for basic auth copied
// from net/http.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Base sets the rawURL. If you intend to extend the url with Path,
// baseUrl should be specified with a trailing slash.
func (s *Client) Base(rawURL ...string) *Client {
	if slices.GreaterThen(0, rawURL) {
		s.rawURL = slices.First(rawURL...)
	}

	return s
}

// Path extends the rawURL with the given path by resolving the reference to
// an absolute URL. If parsing errors occur, the rawURL is left unmodified.
func (s *Client) Path(path string) *Client {
	baseURL, baseErr := url.Parse(s.rawURL)
	pathURL, pathErr := url.Parse(path)

	if baseErr == nil && pathErr == nil {
		s.rawURL = baseURL.ResolveReference(pathURL).String()
		return s
	}

	return s
}

// QueryStruct appends the given queryStruct to the Client's queryStructs.
// The queryStruct argument should be a pointer to a url tagged struct. See
// https://godoc.org/github.com/google/go-querystring/query for details. Any error
// encoding the queryStruct will be returned when creating a request (see
// Request()).
func (s *Client) QueryStruct(queryStruct any) *Client {
	if queryStruct != nil {
		s.queryStructs = append(s.queryStructs, queryStruct)
	}
	return s
}

// Body sets the Client's body. The body value will be set as the Body on new
// requests (see Request()).
// If the provided body is also an io.Closer, the request Body will be closed
// by http.Client methods.
func (s *Client) Body(body io.Reader) *Client {
	if body == nil {
		return s
	}

	return s.BodyProvider(bodyProvider{body: body})
}

// BodyProvider sets the Client's body provider.
func (s *Client) BodyProvider(body BodyProvider) *Client {
	if body == nil {
		return s
	}
	s.bodyProvider = body

	ct := body.ContentType()
	if ct != "" {
		s.Set(contentType, ct)
	}

	return s
}

// BodyJSON sets the Client's bodyJSON. The value pointed to by the bodyJSON
// will be JSON encoded as the Body on new requests (see Request()).
// The bodyJSON argument should be a pointer to a JSON tagged struct. See
// https://golang.org/pkg/encoding/json/#MarshalIndent for details.
func (s *Client) BodyJSON(bodyJSON any) *Client {
	if bodyJSON == nil {
		return s
	}

	return s.BodyProvider(jsonBodyProvider{payload: bodyJSON})
}

// BodyForm sets the Client's bodyForm. The value pointed to by the bodyForm
// will be url encoded as the Body on new requests (see Request()).
// The bodyForm argument should be a pointer to a url tagged struct. See
// https://godoc.org/github.com/google/go-querystring/query for details.
func (s *Client) BodyForm(bodyForm any) *Client {
	if bodyForm == nil {
		return s
	}

	return s.BodyProvider(formBodyProvider{payload: bodyForm})
}

// Request returns a new http.Request created with the Sling properties.
// Returns any errors parsing the rawURL, encoding query structs, encoding
// the body, or creating the http.Request.
func (s *Client) Request(ctx context.Context) (*http.Request, error) {
	reqURL, err := url.Parse(s.rawURL)
	if err != nil {
		return nil, err
	}

	err = addQueryStructs(reqURL, s.queryStructs)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if s.bodyProvider != nil {
		body, err = s.bodyProvider.Body()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, s.method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}

	addHeaders(req, s.header)

	return req, err
}

// addQueryStructs parses url tagged query structs using go-querystring to
// encode them to url.Values and format them onto the url.RawQuery. Any
// query parsing or encoding errors are returned.
func addQueryStructs(reqURL *url.URL, queryStructs []any) error {
	urlValues, err := url.ParseQuery(reqURL.RawQuery)
	if err != nil {
		return err
	}
	// encodes query structs into a url.Values map and merges maps
	for _, queryStruct := range queryStructs {
		queryValues, err := goquery.Values(queryStruct)
		if err != nil {
			return err
		}
		for key, values := range queryValues {
			for _, value := range values {
				urlValues.Add(key, value)
			}
		}
	}
	// url.Values format to a sorted "url encoded" string, e.g. "key=val&foo=bar"
	reqURL.RawQuery = urlValues.Encode()
	return nil
}

// addHeaders adds the key, value pairs from the given http.Header to the
// request. Values for existing keys are appended to the keys values.
func addHeaders(req *http.Request, header http.Header) {
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

// CompletionDecoder sets the Client's completion decoder.
func (s *Client) CompletionDecoder(decoder CompletionDecoder) *Client {
	if decoder == nil {
		return s
	}

	s.completionDecoder = decoder

	return s
}

// CompletionChunksDecoder sets the Client's completion chunks decoder.
func (s *Client) CompletionChunksDecoder(decoder CompletionChunksDecoder) *Client {
	if decoder == nil {
		return s
	}

	s.completionChunksDecoder = decoder

	return s
}

// ResponseDecoder sets the Client's response decoder.
func (s *Client) ResponseDecoder(decoder ResponseDecoder) *Client {
	if decoder == nil {
		return s
	}

	s.responseDecoder = decoder

	return s
}

// Complete is completing the request with the decoder set by completionDecoder.
func (s *Client) Complete(ctx context.Context, prompt *Prompt) (*Completion, error) {
	req, err := s.BodyJSON(prompt).Request(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	completion, err := s.completionDecoder.Decode(resp)
	if err != nil {
		return nil, err
	}

	return completion, nil
}

// CompleteChunked is completing the request with the decoder set by completionDecoder.
// It returns a channel that will receive the completion chunks as they are received.
func (s *Client) CompleteChunked(ctx context.Context, prompt *Prompt) iter.Seq2[*CompletionChunk, error] {
	return func(yield func(*CompletionChunk, error) bool) {
		req, err := s.BodyJSON(prompt).Request(ctx)
		if err != nil {
			yield(nil, err)
			return
		}

		resp, err := s.httpClient.Do(req)
		if err != nil {
			yield(nil, err)
			return
		}
		defer resp.Body.Close()

		for chunk, err := range s.completionChunksDecoder.Decode(resp) {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(chunk, nil) {
				return
			}
		}
	}
}

// ReceiveSuccess creates a new HTTP request and returns the response. Success
// responses (2XX) are JSON decoded into the value pointed to by successV.
// Any error creating the request, sending it, or decoding a 2XX response
// is returned.
func (s *Client) ReceiveSuccess(ctx context.Context, successV any) (*http.Response, error) {
	return s.Receive(ctx, successV, nil)
}

// Receive creates a new HTTP request and returns the response. Success
// responses (2XX) are JSON decoded into the value pointed to by successV and
// other responses are JSON decoded into the value pointed to by failureV.
// If the status code of response is 204(no content) or the Content-Length is 0,
// decoding is skipped. Any error creating the request, sending it, or decoding
// the response is returned.
// Receive is shorthand for calling Request and Do.
func (s *Client) Receive(ctx context.Context, successV, failureV any) (*http.Response, error) {
	req, err := s.Request(ctx)
	if err != nil {
		return nil, err
	}

	return s.Do(req, successV, failureV)
}

// Do sends an HTTP request and returns the response. Success responses (2XX)
// are JSON decoded into the value pointed to by successV and other responses
// are JSON decoded into the value pointed to by failureV.
// If the status code of response is 204(no content) or the Content-Length is 0,
// decoding is skipped. Any error sending the request or decoding the response
// is returned.
func (s *Client) Do(req *http.Request, successV, failureV any) (*http.Response, error) {
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	// when err is nil, resp contains a non-nil resp.Body which must be closed
	defer resp.Body.Close()

	// The default HTTP client's Transport may not
	// reuse HTTP/1.x "keep-alive" TCP connections if the Body is
	// not read to completion and closed.
	// See: https://golang.org/pkg/net/http/#Response
	defer io.Copy(io.Discard, resp.Body)

	// Don't try to decode on 204s or Content-Length is 0
	if resp.StatusCode == http.StatusNoContent || resp.ContentLength == 0 {
		return resp, nil
	}

	// Decode from json
	if successV != nil || failureV != nil {
		err = decodeResponse(resp, s.responseDecoder, successV, failureV)
	}

	return resp, err
}

// decodeResponse decodes response Body into the value pointed to by successV
// if the response is a success (2XX) or into the value pointed to by failureV
// otherwise. If the successV or failureV argument to decode into is nil,
// decoding is skipped.
// Caller is responsible for closing the resp.Body.
func decodeResponse(resp *http.Response, decoder ResponseDecoder, successV, failureV any) error {
	if code := resp.StatusCode; 200 <= code && code <= 299 {
		if successV != nil {
			return decoder.Decode(resp, successV)
		}
	} else {
		if failureV != nil {
			return decoder.Decode(resp, failureV)
		}
	}

	return nil
}
