package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syndis-software/aftra-cli/pkg/openapi"
)

func Test_ExecuteLog_Single(t *testing.T) {
	type test struct {
		serverResponse        int
		serverResponseContent string
		expectedOutput        string
		errorExpected         bool
	}

	tests := map[string]test{
		"success": {serverResponse: 200, serverResponseContent: "", expectedOutput: "", errorExpected: false},
		"401":     {serverResponse: 401, serverResponseContent: "", expectedOutput: "Error: unauthorized\n", errorExpected: true},
		"403":     {serverResponse: 403, serverResponseContent: "", expectedOutput: "Error: forbidden\n", errorExpected: true},
		"422": {
			serverResponse:        422,
			serverResponseContent: "{\"detail\":[{\"loc\":[\"body\",0,\"messages\"],\"msg\":\"field required\",\"type\":\"value_error.missing\"}]}",
			expectedOutput:        "Error: validation error: [{\"loc\":[\"body\",0,\"messages\"],\"msg\":\"field required\",\"type\":\"value_error.missing\"}]\n",
			errorExpected:         true,
		},
		"500": {serverResponse: 500, serverResponseContent: "", expectedOutput: "Error: server error: 500\n", errorExpected: true},
	}

	header := make(http.Header, 1)
	header.Set("Content-Type", "application/json")

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := make(map[string][]Response)
			mockDoer := &MockHTTP{Responses: m}
			// Company is blank in tests
			mockDoer.AddResponse("/api/integrations/syndis-scan/scan-name/logs", Response{
				Response: http.Response{
					StatusCode: tc.serverResponse,
					Status:     "",
					Body:       ioutil.NopCloser(bytes.NewBufferString(tc.serverResponseContent)),
					Header:     header,
				},
				ResponseError: nil,
			})

			actual := new(bytes.Buffer)
			rootCmd.SetOut(actual)
			rootCmd.SetErr(actual)
			rootCmd.SetArgs([]string{"log", "syndis", "scan-name", "My log message"})

			ctx := context.WithValue(context.Background(), doerKey, mockDoer)
			logCmd.SetContext(ctx)

			err := rootCmd.ExecuteContext(ctx)

			if tc.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}

			assert.Equal(t, len(mockDoer.Requests), 1)

			body, _ := ioutil.ReadAll(mockDoer.Requests[0].Body)
			var submitted []openapi.SubmitLogEvent
			_ = json.Unmarshal(body, &submitted)

			assert.Equal(t, len(submitted), 1)
			assert.Equal(t, "My log message", submitted[0].Message)
			assert.Equal(t, tc.expectedOutput, actual.String())
		})
	}

}

func Test_ExecuteLog_Stdin(t *testing.T) {

	type test struct {
		serverResponse        int
		serverResponseContent string
		expectedErrOutput     string
	}

	tests := map[string]test{
		"success": {serverResponse: 200, serverResponseContent: "", expectedErrOutput: ""},
		"401":     {serverResponse: 403, serverResponseContent: "", expectedErrOutput: "Error: forbidden\n"},
		"422": {
			serverResponse:        422,
			serverResponseContent: "{\"detail\":[{\"loc\":[\"body\",0,\"messages\"],\"msg\":\"field required\",\"type\":\"value_error.missing\"}]}",
			expectedErrOutput:     "Error: validation error: [{\"loc\":[\"body\",0,\"messages\"],\"msg\":\"field required\",\"type\":\"value_error.missing\"}]\n",
		},
		"500": {serverResponse: 500, serverResponseContent: "", expectedErrOutput: "Error: server error: 500\n"},
	}

	header := make(http.Header, 1)
	header.Set("Content-Type", "application/json")

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			stdinInput := strings.NewReader("abcde\nfoobar")

			m := make(map[string][]Response)
			mockDoer := &MockHTTP{Responses: m}
			// Company is blank in tests
			mockDoer.AddResponse("/api/integrations/syndis-scan/scan-name/logs", Response{
				Response: http.Response{
					StatusCode: tc.serverResponse,
					Status:     "",
					Body:       ioutil.NopCloser(bytes.NewBufferString(tc.serverResponseContent)),
					Header:     header,
				},
				ResponseError: nil,
			})

			outStd := new(bytes.Buffer)
			outErr := new(bytes.Buffer)
			rootCmd.SetOut(outStd)
			rootCmd.SetErr(outErr)
			rootCmd.SetIn(stdinInput)
			rootCmd.SetArgs([]string{"log", "syndis", "scan-name"})

			ctx := context.WithValue(context.Background(), doerKey, mockDoer)
			logCmd.SetContext(ctx)

			err := rootCmd.ExecuteContext(ctx)

			assert.Equal(t, nil, err)
			assert.Equal(t, len(mockDoer.Requests), 1)

			body, _ := ioutil.ReadAll(mockDoer.Requests[0].Body)
			var submitted []openapi.SubmitLogEvent
			_ = json.Unmarshal(body, &submitted)

			assert.Equal(t, len(submitted), 2)
			assert.Equal(t, "abcde", submitted[0].Message)
			assert.Equal(t, "foobar", submitted[1].Message)
			assert.Equal(t, "", outStd.String())
			assert.Contains(t, outErr.String(), tc.expectedErrOutput)
		})
	}
}
