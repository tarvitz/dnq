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

// Call performs an generic http rest api calls.
func (api *Client) Call(request *APIRequest, in interface{}) (err error) {
	var (
		response    *http.Response
		httpRequest *http.Request
	)

	url := api.URL("/%s", request.APIMethod)
	httpRequest, _ = http.NewRequest(request.Method, url, request.Body)
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
func (api *Client) Retrieve(endpoint string, in interface{}) (err error) {
	request := &APIRequest{Method: "GET", APIMethod: MethodType(endpoint)}
	return api.Call(request, &in)
}

// Upload does a multipart/form-data call to the given URL.
func (api *Client) Upload(method MethodType, values map[string]io.Reader) (
	message *Message, err error) {

	var (
		body        io.Reader
		contentType string
	)
	body, contentType = makeMultipartFormDataPayload(values)

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
	body io.Reader, contentType string) {
	var closeStack []interface{}
	defer safeClose(closeStack)()

	buffer := bytes.NewBuffer([]byte{})
	writer := multipart.NewWriter(buffer)

	for field, reader := range values {
		var fieldWriter io.Writer
		closeStack = append(closeStack, fieldWriter)

		switch in := reader.(type) {
		case *os.File:
			// ignoring errors due to very hard and complicated way of their
			// producing. Note, there still can be an issue, but on the present moment
			// (an initial version) the risks are counted as insignificant.
			fieldWriter, _ = writer.CreateFormFile(field, in.Name())
		default:
			fieldWriter, _ = writer.CreateFormField(field)
		}
		_, _ = io.Copy(fieldWriter, reader)
	}
	_ = writer.Close()
	return buffer, writer.FormDataContentType()
}
