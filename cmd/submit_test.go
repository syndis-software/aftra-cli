package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syndis-software/aftra-cli/pkg/openapi"
)

func Test_ExecuteSubmit_Single_ServerResponseHandling(t *testing.T) {
	type test struct {
		serverResponse        int
		serverResponseContent string
		expectedOutput        string
		errorExpected         bool
	}

	tests := map[string]test{
		"success": {
			serverResponse:        200,
			serverResponseContent: "",
			expectedOutput:        "",
			errorExpected:         false,
		},
		"401": {
			serverResponse:        401,
			serverResponseContent: "",
			expectedOutput:        "Error: unauthorized\n",
			errorExpected:         true},
		"403": {
			serverResponse:        403,
			serverResponseContent: "",
			expectedOutput:        "Error: forbidden\n",
			errorExpected:         true},
		"422": {
			serverResponse:        422,
			serverResponseContent: "{\"detail\":[{\"loc\":[\"body\",0,\"messages\"],\"msg\":\"field required\",\"type\":\"value_error.missing\"}]}",
			expectedOutput:        "Error: validation error: [{\"loc\":[\"body\",0,\"messages\"],\"msg\":\"field required\",\"type\":\"value_error.missing\"}]\n",
			errorExpected:         true,
		},
		"500": {
			serverResponse:        500,
			serverResponseContent: "",
			expectedOutput:        "Error: server error: 500\n",
			errorExpected:         true},
	}

	header := make(http.Header, 1)
	header.Set("Content-Type", "application/json")
	submitCmd_filename = ""
	submitCmd_message = ""
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := make(map[string][]Response)
			mockDoer := &MockHTTP{Responses: m}
			mockDoer.AddResponse("/api/integrations/syndis-scan/scan-name/scan", Response{
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
			rootCmd.SetArgs([]string{"submit", "syndis", "scan-name", "--message", `[{"folder_name":"Islenskur_Texti","scan_id":"117","scan_uuid":"d776c35d-7d0e-b699-d035-a0120d2dfe28a871ce19964f5207","scan_start":"2023-04-04 08:25:22","scan_end":"2023-04-04 08:32:35","file_name":"External.csv","plugin_id":10114,"cve":"CVE-1999-0524","cvss":0,"risk":"None","host":"153.92.147.3","protocol":"icmp","port":"0","name":"ICMP Timestamp Request Remote Date Disclosure","synopsis":"It is possible to determine the exact time set on the remote host.","description":"The remote host answers to an ICMP timestamp request.  This allows an\nattacker to know the date that is set on the targeted machine, which\nmay assist an unauthenticated, remote attacker in defeating time-based\nauthentication protocols.\n\nTimestamps returned from machines running Windows Vista / 7 / 2008 /\n2008 R2 are deliberately incorrect, but usually within 1000 seconds of\nthe actual system time.","solution":"Filter out the ICMP timestamp requests (13), and the outgoing ICMP\\ntimestamp replies (14).","see_also":"","plugin_output":"The remote clock is synchronized with the local clock.\\n","stig_severity":"","cvss_v3_0_base_score":0,"cvss_temporal_score":0,"cvss_v3_0_temporal_score":0,"risk_factor":"None","bid":"","xref":"CWE:200","mskb":"","plugin_publication_date":"1999/08/01","plugin_modification_date":"2019/10/04","metasploit":"","core_impact":"","canvas":""}]`})

			ctx := context.WithValue(context.Background(), doerKey, mockDoer)
			submitCmd.SetContext(ctx)

			err := rootCmd.ExecuteContext(ctx)

			if tc.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}

			assert.Equal(t, len(mockDoer.Requests), 1)

			body, _ := ioutil.ReadAll(mockDoer.Requests[0].Body)

			var submitted openapi.BodySubmitScanResults
			_ = json.Unmarshal(body, &submitted)

			events, _ := submitted.Events.AsBodySubmitScanResultsEvents1()
			assert.Equal(t, 1, len(events))
			assert.Equal(t, tc.expectedOutput, actual.String())
		})
	}

}

func Test_ExecuteSubmit_Single_JsonParsing(t *testing.T) {
	type test struct {
		message        string
		expectedOutput string
		errorExpected  bool
	}

	tests := map[string]test{
		"success": {
			message:        "[{}]",
			expectedOutput: "",
			errorExpected:  false,
		},
		"malformed-json": {
			message:        "ooo",
			expectedOutput: "Error: invalid character 'o' looking for beginning of value\n",
			errorExpected:  true,
		},
	}

	header := make(http.Header, 1)
	header.Set("Content-Type", "application/json")
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			m := make(map[string][]Response)
			mockDoer := &MockHTTP{Responses: m}
			mockDoer.AddResponse("/api/integrations/syndis-scan/scan-name/scan", Response{
				Response: http.Response{
					StatusCode: 200,
					Status:     "",
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
					Header:     header,
				},
				ResponseError: nil,
			})

			actual := new(bytes.Buffer)
			rootCmd.SetOut(actual)
			rootCmd.SetErr(actual)
			rootCmd.SetArgs([]string{"submit", "syndis", "scan-name", "--message", tc.message})

			ctx := context.WithValue(context.Background(), doerKey, mockDoer)
			submitCmd.SetContext(ctx)

			err := rootCmd.ExecuteContext(ctx)

			var expectedRequests int

			if tc.errorExpected {
				assert.NotNil(t, err)
				expectedRequests = 0
			} else {
				assert.Equal(t, nil, err)
				expectedRequests = 1
			}

			assert.Equal(t, expectedRequests, len(mockDoer.Requests))
			assert.Equal(t, tc.expectedOutput, actual.String())
		})
	}

}

func Test_ExecuteSubmit_File(t *testing.T) {
	type test struct {
		fileContent    string
		expectedOutput string
		errorExpected  bool
	}

	tests := map[string]test{
		"success": {
			fileContent:    "[{}]",
			expectedOutput: "",
			errorExpected:  false,
		},
		"malformed-json": {
			fileContent:    "ooo",
			expectedOutput: "Error: invalid character 'o' looking for beginning of value\n",
			errorExpected:  true,
		},
	}

	submitCmd_filename = ""
	submitCmd_message = ""
	header := make(http.Header, 1)
	header.Set("Content-Type", "application/json")

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			file, err := ioutil.TempFile("/tmp", "submit-test")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(file.Name())
			file.Write([]byte(tc.fileContent))

			m := make(map[string][]Response)
			mockDoer := &MockHTTP{Responses: m}
			mockDoer.AddResponse("/api/integrations/syndis-scan/scan-name/scan", Response{
				Response: http.Response{
					StatusCode: 200,
					Status:     "",
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
					Header:     header,
				},
				ResponseError: nil,
			})
			upload_response := &openapi.BlobSignedUploadURLResponse{
				Bucket: "bucket",
				Key:    "key",
				Url:    "http://foo.com",
			}
			b, _ := json.Marshal(upload_response)
			mockDoer.AddResponse("/api/companies//blobs/upload", Response{
				Response: http.Response{
					StatusCode: 200,
					Status:     "",
					Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
					Header:     header,
				},
				ResponseError: nil,
			})

			actual := new(bytes.Buffer)
			rootCmd.SetOut(actual)
			rootCmd.SetErr(actual)
			rootCmd.SetArgs([]string{"submit", "syndis", "scan-name", "--filename", file.Name()})

			ctx := context.WithValue(context.Background(), doerKey, mockDoer)
			submitCmd.SetContext(ctx)

			err = rootCmd.ExecuteContext(ctx)

			var expectedRequests int

			if tc.errorExpected {
				assert.NotNil(t, err)
				expectedRequests = 0
			} else {
				assert.Equal(t, nil, err)
				expectedRequests = 1
			}

			assert.Equal(t, expectedRequests, mockDoer.CountRequests("/api/integrations/syndis-scan/scan-name/scan"))
			assert.Equal(t, tc.expectedOutput, actual.String())
		})
	}

}
