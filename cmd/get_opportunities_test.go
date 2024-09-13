package cmd

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteGetOpportunities(t *testing.T) {
	type test struct {
		serverResponse       int
		limitArg             string
		serverResponseData   []string
		expectedOutput       string
		expectedRequestCount int
	}

	tests := []test{
		{
			serverResponse:       200,
			limitArg:             "100",
			serverResponseData:   []string{`{"opportunities": [{"a": 1}, {"b": 2}], "total": 2}`},
			expectedOutput:       "{\"a\":1}\n{\"b\":2}\n",
			expectedRequestCount: 1,
		},
		{
			serverResponse:       200,
			limitArg:             "2",
			serverResponseData:   []string{`{"opportunities": [{"a": 1}], "total": 2}`, `{"opportunities": [{"b": 2}], "total": 2}`},
			expectedOutput:       "{\"a\":1}\n{\"b\":2}\n",
			expectedRequestCount: 2,
		},
	}

	for _, tc := range tests {
		t.Run(name, func(t *testing.T) {
			header := make(http.Header, 1)
			header.Set("Content-Type", "application/json")

			m := make(map[string][]Response)
			mockDoer := &MockHTTP{Responses: m}
			for _, responseData := range tc.serverResponseData {
				mockDoer.AddResponse("/api/companies//opportunities/v3", Response{
					Response: http.Response{
						StatusCode: tc.serverResponse,
						Status:     "200",
						Body:       ioutil.NopCloser(bytes.NewBufferString(responseData)),
						Header:     header,
					},
					ResponseError: nil,
				})
			}

			actual := new(bytes.Buffer)
			rootCmd.SetOut(actual)
			rootCmd.SetErr(actual)
			rootCmd.SetArgs([]string{"get", "opportunities", "--updated-since", "2024-01-01T00:00:00Z", "--limit", tc.limitArg})

			ctx := context.WithValue(context.Background(), doerKey, mockDoer)
			getOpportunitiesCmd.SetContext(ctx)

			err := rootCmd.ExecuteContext(ctx)

			assert.Equal(t, nil, err)
			assert.Equal(t, tc.expectedRequestCount, len(mockDoer.Requests))
			assert.Equal(t, tc.expectedOutput, actual.String())
		})
	}

}

func Test_ValidationGetOpportunities(t *testing.T) {
	type test struct {
		updatedSinceArg       string
		limitArg              string
		expectedOutputPartial string
	}

	tests := []test{
		{
			updatedSinceArg:       "wrong-format",
			limitArg:              "100",
			expectedOutputPartial: "Error: invalid time format",
		},
		{
			updatedSinceArg:       "2024-01-01T00:00:00Z",
			limitArg:              "abc",
			expectedOutputPartial: "Error: invalid argument \"abc\" for \"--limit\" flag",
		},
		{
			updatedSinceArg:       "2024-01-01T00:00:00Z",
			limitArg:              "-4",
			expectedOutputPartial: "Error: limit should be -1 (everything) or less than 1000: -4",
		},
		{
			updatedSinceArg:       "2024-01-01T00:00:00Z",
			limitArg:              "1003",
			expectedOutputPartial: "Error: limit should be -1 (everything) or less than 1000: 1003",
		},
	}

	for _, tc := range tests {

		mockDoer := &MockHTTP{}
		actual := new(bytes.Buffer)
		rootCmd.SetOut(actual)
		rootCmd.SetErr(actual)
		rootCmd.SetArgs([]string{"get", "opportunities", "--updated-since", tc.updatedSinceArg, "--limit", tc.limitArg})

		ctx := context.WithValue(context.Background(), doerKey, mockDoer)
		getOpportunitiesCmd.SetContext(ctx)

		rootCmd.ExecuteContext(ctx)

		assert.Contains(t, actual.String(), tc.expectedOutputPartial)
	}

}
