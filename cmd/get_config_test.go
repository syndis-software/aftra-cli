package cmd

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteGetTokenConfig(t *testing.T) {
	type test struct {
		serverResponse     int
		serverResponseData string
		expectedOutput     string
	}

	tests := []test{
		{
			serverResponse:     200,
			serverResponseData: `{"ranges":"0/1","type":"INTERNAL"}`,
			expectedOutput:     `{"ranges":"0/1","type":"INTERNAL"}`,
		},
	}

	for _, tc := range tests {
		header := make(http.Header, 1)
		header.Set("Content-Type", "application/json")

		m := make(map[string][]Response)
		mockDoer := &MockHTTP{Responses: m}

		mockDoer.AddResponse("/api/integrations/syndis-scan/scan-name/config", Response{
			Response: http.Response{
				StatusCode: tc.serverResponse,
				Status:     "",
				Body:       ioutil.NopCloser(bytes.NewBufferString(tc.serverResponseData)),
				Header:     header,
			},
			ResponseError: nil,
		})

		actual := new(bytes.Buffer)
		rootCmd.SetOut(actual)
		rootCmd.SetErr(actual)
		rootCmd.SetArgs([]string{"get", "config", "syndis", "scan-name"})

		ctx := context.WithValue(context.Background(), doerKey, mockDoer)
		getConfigCmd.SetContext(ctx)

		err := rootCmd.ExecuteContext(ctx)

		assert.Equal(t, nil, err)
		assert.Equal(t, len(mockDoer.Requests), 1)
		assert.Equal(t, tc.expectedOutput, actual.String())
	}

}

func Test_ExecuteGetAllConfigsForType(t *testing.T) {
	type test struct {
		serverResponse     int
		serverResponseData string
		expectedOutput     string
	}

	tests := []test{
		{
			serverResponse: 200,
			serverResponseData: `{
				"entities": [
				  {
					"pk": "Company-90fcc7d5-52f3-4c90-b6e2-c2319c027ae4",
					"sk": "SyndisScan-32f181df-4095-4d69-b522-a5b8bee75ff4",
					"entityType": "SyndisScan",
					"created": "2023-05-03T12:41:26.026574",
					"updated": "2023-05-03T12:41:26.026648",
					"name": "New1Jakob3",
					"config": {
					  "type": "INTERNAL",
					  "ranges": "127.0.0.1"
					}
				  },
				  {
					"pk": "Company-90fcc7d5-52f3-4c90-b6e2-c2319c027ae4",
					"sk": "SyndisScan-a0a8e2e7-2df7-4161-9d88-05cfd436cf97",
					"entityType": "SyndisScan",
					"created": "2023-05-03T12:40:57.704505",
					"updated": "2023-05-03T12:40:57.704582",
					"name": "New1Jakob2",
					"config": {
					  "type": "INTERNAL",
					  "ranges": "127.0.0.1"
					}
				  }
				]
			}`,

			expectedOutput: `{"config":{"ranges":"127.0.0.1","type":"INTERNAL"},"name":"New1Jakob3"}
{"config":{"ranges":"127.0.0.1","type":"INTERNAL"},"name":"New1Jakob2"}
`,
		},
	}

	for _, tc := range tests {
		header := make(http.Header, 1)
		header.Set("Content-Type", "application/json")

		m := make(map[string][]Response)
		mockDoer := &MockHTTP{Responses: m}
		// Company is blank in tests
		mockDoer.AddResponse("/api/companies//syndis-scans", Response{
			Response: http.Response{
				StatusCode: tc.serverResponse,
				Status:     "",
				Body:       ioutil.NopCloser(bytes.NewBufferString(tc.serverResponseData)),
				Header:     header,
			},
			ResponseError: nil,
		})

		actual := new(bytes.Buffer)
		rootCmd.SetOut(actual)
		rootCmd.SetErr(actual)
		rootCmd.SetArgs([]string{"get", "config", "syndis"})

		ctx := context.WithValue(context.Background(), doerKey, mockDoer)
		getConfigCmd.SetContext(ctx)

		err := rootCmd.ExecuteContext(ctx)

		assert.Equal(t, nil, err)
		assert.Equal(t, len(mockDoer.Requests), 1)
		assert.Equal(t, tc.expectedOutput, actual.String())
	}

}
