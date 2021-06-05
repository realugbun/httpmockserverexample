package httpmockserver

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	// Example of using this package. It has its own assert for testing but it can be replaced by other packages without affecting the mock server
	"git.sr.ht/~ewintr/go-kit/test"
)

// This is a full working example based on the blog post: https://erikwinter.nl/articles/2020/unit-test-outbound-http-requests-in-golang/

// See the post for more details

// Note: The mock server doesn't validate requests. It sends back
// a string in the body according to the test. You will need to make
// your own tests for validating the requests using the record variable
// below

func TestFooClientDoStuffJSON(t *testing.T) {

	// Setup the test scenarios
	path := "/"
	username = "username"
	password = "password"

	for _, tc := range []struct {
		name      string
		param     string
		respCode  int
		respBody  string
		expErr    error
		expResult string
	}{
		{
			name:     "upstream failure",
			respCode: http.StatusInternalServerError,
			expErr:   ErrUpstreamFailure,
		},
		{
			name:      "valid response to bar", // Name of the test
			param:     "bar",                   // The param that will be sent in the request body ie {"param":"bar"}
			respCode:  http.StatusOK,           // The status code the mock server will send back
			respBody:  `{"result":"ok"}`,       // The body the mock server will send back
			expResult: "ok",                    // The expected output from the DoStuff method
			expErr:    nil,                     // The expected error from the method.
		},
		{
			name:      "valid response to baz",
			param:     "baz",
			respCode:  http.StatusOK,
			respBody:  `{"result":"also ok"}`,
			expResult: "also ok",
			expErr:    nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Record will record the requests sent to the mock server
			var record test.MockAssertion

			// Setup the mock server
			mockServer := test.NewMockServer(&record, test.MockServerProcedure{
				URI:        path,
				HTTPMethod: http.MethodPost,
				Response: test.MockResponse{
					StatusCode: tc.respCode,
					Body:       []byte(tc.respBody),
				},
			})

			// Make a new client using our program
			client := NewFooClient(mockServer.URL, username, password)

			// Make the call to the server
			actResult, actErr := client.DoStuffJSON(tc.param)

			// check the request was done
			test.Equals(t, 1, record.Hits(path, http.MethodPost))

			// check request body
			expBody := fmt.Sprintf(`{"param":%q}`, tc.param)
			actBody := string(record.Body(path, http.MethodPost)[0])
			test.Equals(t, expBody, actBody)

			// check request headers
			expHeaders := []http.Header{{
				"Authorization": []string{"Basic dXNlcm5hbWU6cGFzc3dvcmQ="},
				"Content-Type":  []string{"application/json;charset=utf-8"},
			}}
			test.Equals(t, expHeaders, record.Headers(path, http.MethodPost))

			// check what the method returns
			test.Equals(t, true, errors.Is(actErr, tc.expErr))
			test.Equals(t, tc.expResult, actResult)
		})
	}
}

// Shows how the mock server can be used for APIs that do not use JSON
func TestFooClientDoStuffXML(t *testing.T) {
	path := "/"
	username = "username"
	password = "password"

	for _, tc := range []struct {
		name      string
		param     string
		respCode  int
		respBody  string
		expErr    error
		expResult string
	}{
		{
			name:     "upstream failure",
			respCode: http.StatusInternalServerError,
			expErr:   ErrUpstreamFailure,
		},
		{
			name:      "valid response to bar",
			param:     "bar",
			respCode:  http.StatusOK,
			respBody:  `<xml><body><Result>ok</Result></body></xml>`,
			expResult: "ok",
		},
		{
			name:      "valid response to baz",
			param:     "baz",
			respCode:  http.StatusOK,
			respBody:  `<xml><body><Result>also ok</Result></body></xml>`,
			expResult: "also ok",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var record test.MockAssertion
			mockServer := test.NewMockServer(&record, test.MockServerProcedure{
				URI:        path,
				HTTPMethod: http.MethodPost,
				Response: test.MockResponse{
					StatusCode: tc.respCode,
					Body:       []byte(tc.respBody),
				},
			})

			client := NewFooClient(mockServer.URL, username, password)

			actResult, actErr := client.DoStuffXML(tc.param)

			// check request was done
			test.Equals(t, 1, record.Hits(path, http.MethodPost))

			// check request body
			expBody := fmt.Sprintf(`<Request>%s</Request>`, tc.param)
			actBody := string(record.Body(path, http.MethodPost)[0])
			test.Equals(t, expBody, actBody)

			// check request headers
			expHeaders := []http.Header{{
				"Authorization": []string{"Basic dXNlcm5hbWU6cGFzc3dvcmQ="},
				"Content-Type":  []string{contentTypeXML},
			}}
			test.Equals(t, expHeaders, record.Headers(path, http.MethodPost))

			// check result
			test.Equals(t, true, errors.Is(actErr, tc.expErr))
			test.Equals(t, tc.expResult, actResult)
		})
	}
}
