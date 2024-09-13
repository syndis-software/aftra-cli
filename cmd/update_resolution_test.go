package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/syndis-software/aftra-cli/pkg/openapi"
)

func Test_ExecuteUpdateResolution(t *testing.T) {
	type test struct {
		serverResponse int
		expectedOutput string
		uid            string
	}

	tests := []test{
		{
			uid:            "123",
			serverResponse: 200,
			expectedOutput: "",
		},
	}

	for _, tc := range tests {
		t.Run(name, func(to *testing.T) {
			header := make(http.Header, 1)
			header.Set("Content-Type", "application/json")

			m := make(map[string][]Response)
			mockDoer := &MockHTTP{Responses: m}
			mockDoer.AddResponse(fmt.Sprintf("/api/companies//opportunities/%s/", tc.uid), Response{
				Response: http.Response{
					StatusCode: tc.serverResponse,
					Status:     fmt.Sprint(tc.serverResponse),
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
					Header:     header,
				},
				ResponseError: nil,
			})
			actual := new(bytes.Buffer)
			rootCmd.SetOut(actual)
			rootCmd.SetErr(actual)
			rootCmd.SetArgs([]string{"update", "resolution", tc.uid, "accepted_risk", "--comment", "Worked all night", "--due-date", "2024-02-04"})

			ctx := context.WithValue(context.Background(), doerKey, mockDoer)
			updateResolutionsCmd.SetContext(ctx)

			err := rootCmd.ExecuteContext(ctx)

			assert.Equal(to, nil, err)
			assert.Equal(to, tc.expectedOutput, actual.String())
			assert.Equal(to, 1, mockDoer.CountRequests(fmt.Sprintf("/api/companies//opportunities/%s/", tc.uid)))
		})
	}

}

func Test_UpdateResolutionValidation(t *testing.T) {
	type test struct {
		resolution    string
		expectedValid bool
		// Ignored if expectedValid is false
		expectedValue openapi.OpportunityResolution
	}

	tests := []test{
		{
			resolution:    "resolved",
			expectedValue: openapi.Resolved,
			expectedValid: true,
		},
		{
			resolution:    "accepted_risk",
			expectedValue: openapi.AcceptedRisk,
			expectedValid: true,
		},
		{
			resolution:    "false_positive",
			expectedValue: openapi.FalsePositive,
			expectedValid: true,
		},
		{
			resolution:    "unacknowledged",
			expectedValue: openapi.Unacknowledged,
			expectedValid: true,
		},
		{
			resolution:    "Not there",
			expectedValue: openapi.Resolved,
			expectedValid: false,
		},
		{
			resolution:    "RESOLVED",
			expectedValue: openapi.Resolved,
			expectedValid: false,
		},
	}

	for _, tc := range tests {
		t.Run(name, func(to *testing.T) {
			v, err := validateResolution(tc.resolution)
			if tc.expectedValid {
				assert.Equal(to, nil, err)
				assert.Equal(to, *v, tc.expectedValue)
			} else {
				assert.NotNil(to, err)
			}

		})
	}

}

func Test_UpdateResolutionDueDateValidation(t *testing.T) {
	type test struct {
		dueDateString string
		resolution    openapi.OpportunityResolution
		expectedValue *openapi_types.Date
		expectedValid bool
	}

	jan1st := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []test{
		{
			dueDateString: "2024-01-01",
			resolution:    openapi.Resolved,
			expectedValid: false,
		},
		{
			dueDateString: "2024-01-01",
			resolution:    openapi.FalsePositive,
			expectedValid: false,
		},
		{
			dueDateString: "2024-01-01",
			resolution:    openapi.Unacknowledged,
			expectedValid: false,
		},
		{
			dueDateString: "",
			resolution:    openapi.Resolved,
			expectedValid: true,
			expectedValue: nil,
		},
		{
			dueDateString: "",
			resolution:    openapi.FalsePositive,
			expectedValid: true,
			expectedValue: nil,
		},
		{
			dueDateString: "",
			resolution:    openapi.Unacknowledged,
			expectedValid: true,
			expectedValue: nil,
		},
		{
			dueDateString: "2024-01-01",
			resolution:    openapi.AcceptedRisk,
			expectedValid: true,
			expectedValue: &openapi_types.Date{jan1st},
		},
		{
			dueDateString: "",
			resolution:    openapi.AcceptedRisk,
			expectedValid: false,
		},
		{
			dueDateString: "invalid-date",
			resolution:    openapi.AcceptedRisk,
			expectedValid: false,
		},
	}

	for _, tc := range tests {
		t.Run(name, func(to *testing.T) {
			v, err := validateDueDateWithResolutionCheck(tc.dueDateString, tc.resolution)
			if tc.expectedValid {
				assert.Equal(to, nil, err)
				assert.Equal(to, v, tc.expectedValue)
			} else {
				assert.NotNil(to, err)
			}

		})
	}

}
