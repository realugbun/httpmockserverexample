package httpmockserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

// These are basic functions which exist just to show how to use a mock server for testing

const (
	contentTypeJSON = "application/json;charset=utf-8"
	contentTypeXML  = "application/xml"
)

var ErrUpstreamFailure = errors.New("upstream failure")

var (
	username string
	password string
)

type FooClient struct {
	client  *http.Client
	request *http.Request
}

type ResponseJSONRec struct {
	Result string `json:"result"`
}

type ResponseXMLRec struct {
	XMLName xml.Name `xml:"xml"`
	Text    string   `xml:",chardata"`
	Body    struct {
		Text   string `xml:",chardata"`
		Result string `xml:"Result"`
	} `xml:"body"`
}

type RequestRec struct {
	Param string `json:"param"`
}

func NewFooClient(url, username, password string) *FooClient {

	var newClient FooClient

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return &FooClient{}
	}

	addBasicAuth(request, username, password)

	newClient.client = client
	newClient.request = request

	return &newClient
}

func (fc *FooClient) DoStuffJSON(param string) (string, error) {

	fc.request.Header.Set("Content-Type", contentTypeJSON)

	// Do stuff does something and creates a request body
	body := io.NopCloser(bytes.NewReader([]byte(`{"param":"` + param + `"}`)))

	// Send the request
	fc.request.Body = body
	response, err := fc.client.Do(fc.request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	// Only one example
	if response.StatusCode == http.StatusInternalServerError {
		return "", ErrUpstreamFailure
	}

	// Read the response
	var responseBody ResponseJSONRec
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&responseBody)

	return responseBody.Result, nil
}

func (fc *FooClient) DoStuffXML(param string) (string, error) {

	fc.request.Header.Set("Content-Type", contentTypeXML)

	// Do stuff does something and creates a request body
	body := io.NopCloser(bytes.NewReader([]byte(`<Request>` + param + `</Request>`)))

	// Send the request
	fc.request.Body = body
	response, err := fc.client.Do(fc.request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	// Only one example
	if response.StatusCode == http.StatusInternalServerError {
		return "", ErrUpstreamFailure
	}

	// Read the body
	var responseBody ResponseXMLRec
	decoder := xml.NewDecoder(response.Body)
	decoder.Decode(&responseBody)

	return responseBody.Body.Result, nil
}

func addBasicAuth(r *http.Request, u, p string) {
	auth := u + ":" + p
	c := base64.StdEncoding.EncodeToString([]byte(auth))
	b := "Basic "

	r.Header.Add("Authorization", b+c)
}
