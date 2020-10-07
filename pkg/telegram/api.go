package telegram

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
)

// APIRequest is a plain and simple structure to perform Gitlab DomainAPI requests.
type APIRequest struct {
	Method       string
	APIMethod    MethodType
	ExpectStatus int
	Body         io.Reader
	Query        string
	ContentType  string
}

// Client is a gitlab api client
type Client struct {
	token      string
	httpClient *http.Client

	apiURL string

	// Services
	Inlines *InlineService
}

// WithToken sets token for authentication
func (api *Client) WithToken(token string) *Client {
	api.token = token
	return api
}

// WithURL sets url to send bot api requests to. It's helpful in testing.
func (api *Client) WithURL(url string) *Client {
	api.apiURL = url
	return api
}

// NewClient creates an DomainAPI client.
func NewClient(token string) *Client {
	jar, _ := cookiejar.New(nil)
	client := &Client{
		apiURL:     BotAPIURL,
		token:      token,
		httpClient: &http.Client{Jar: jar},
	}
	client.Inlines = &InlineService{client: client}
	return client
}

// URL returns full REST api endpoint based on apiURL
func (api *Client) URL(format string, args ...interface{}) string {
	rightPart := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s%s", api.apiURL, api.token, rightPart)
}

func (api *Client) request(method, url string, body io.Reader) *http.Request {
	// error has been omitted by an intention due to
	// api client should verify a sanity of api urls before constructing a request.
	request, _ := http.NewRequest(method, url, body)
	return request
}

func (api *Client) updateRequest(request *APIRequest, httpRequest *http.Request) *http.Request {
	if request.Query != "" {
		httpRequest.URL.RawQuery = request.Query
	}
	return httpRequest
}

// CRUD like implementation

// Call performs an generic http rest api calls.
func (api *Client) Call(request *APIRequest, in interface{}) (err error) {
	var response *http.Response

	url := api.URL("/%s", request.APIMethod)
	httpRequest := api.updateRequest(request, api.request(request.Method, url, request.Body))
	httpRequest.Header.Set("Content-Type", orString(request.ContentType, "application/json"))

	response, err = api.httpClient.Do(httpRequest)
	if err != nil {
		return
	}
	defer Close(response.Body)

	content, _ := ioutil.ReadAll(bufio.NewReader(response.Body))

	// use expect-status or fallback to status-ok
	if response.StatusCode != orInt(request.ExpectStatus, http.StatusOK) {
		err = fmt.Errorf("fail to perform [%s][%d]: %s", request.Method, response.StatusCode, content)
		return
	}

	// Go client omits message-body due to:
	// https://tools.ietf.org/html/rfc2616#section-10.2.5
	// see also: https://github.com/golang/go/issues/6685
	// Thus far, even if you will have DELETE methods with a payload, go client won't
	// work with it, so don't even try ;).
	if request.Method != "DELETE" {
		err = json.Unmarshal(content, &in)
	}
	return
}

// Retrieve is a basic GET operation for getting object
// (or in some particular cases, the list of objects)
func (api *Client) Retrieve(endpoint string, body io.Reader, in interface{}) (err error) {
	request := &APIRequest{Method: "GET", APIMethod: MethodType(endpoint), Body: body}
	return api.Call(request, &in)
}

// Upload does a multipart/form-data call to the given URL.
func (api *Client) Upload(method MethodType, values map[string]io.Reader) (
	message *Message, err error) {

	var (
		body        io.Reader
		contentType string
	)
	if body, contentType, err = makeMultipartFormDataPayload(values); err != nil {
		return
	}

	request := &APIRequest{
		Method: "POST", APIMethod: method,
		Body: body, ContentType: contentType,
	}
	response := APIResponse{}
	err = api.Call(request, &response)
	return response.Message, err
}

// ReadUpdate reads an update entry sent via Webhook by telegram
// see also: https://core.telegram.org/bots/api#setwebhook
func ReadUpdate(request *http.Request) (update *Update, err error) {
	var (
		content []byte
	)
	defer Close(request.Body)

	content, err = ioutil.ReadAll(request.Body)
	err = json.Unmarshal(content, &update)
	return
}

func makeMultipartFormDataPayload(values map[string]io.Reader) (
	body io.Reader, contentType string, err error) {
	var closeStack []interface{}
	defer safeClose(closeStack)()

	buffer := bytes.NewBuffer([]byte{})
	writer := multipart.NewWriter(buffer)

	for field, reader := range values {
		var fieldWriter io.Writer
		closeStack = append(closeStack, fieldWriter)

		switch in := reader.(type) {
		case *os.File:
			if fieldWriter, err = writer.CreateFormFile(field, in.Name()); err != nil {
				return
			}
		default:
			if fieldWriter, err = writer.CreateFormField(field); err != nil {
				return
			}
		}
		if _, err = io.Copy(fieldWriter, reader); err != nil {
			return
		}
	}
	_ = writer.Close()
	return buffer, writer.FormDataContentType(), err
}
