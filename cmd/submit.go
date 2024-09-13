/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	openapi "github.com/syndis-software/aftra-cli/pkg/openapi"
)

var (
	submitCmd_filename string
	submitCmd_message  string

	submitCmd = &cobra.Command{
		Use:   "submit [scan-type] [scan-name] [scan-result]",
		Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
		Short: "Submit a json formatted scan result",
		Long: `Submit a json formatted scan result

Submit a scan result in the format for given scan-type. For example in
nessus format for syndis scans.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			client := ctx.Value(clientKey).(*openapi.ClientWithResponses)
			scanType, scanName := args[0], args[1]

			switch {
			case ScanType(scanType) == syndis:

				var results openapi.BodySubmitScanResults
				// Submit a set of scan events
				if submitCmd_filename != "" {
					jsonFile, err := os.Open(submitCmd_filename)
					if err != nil {
						return err
					}
					defer jsonFile.Close()
					contents, _ := ioutil.ReadAll(jsonFile)

					// We attempt to parse the contents to see if its a valid json file
					err = validate_json_file(contents)
					if err != nil {
						return err
					}

					company := ctx.Value(companyKey).(string)
					// Upload the file
					// 1 Get a signed upload url
					uploadInfo, err := openapi.DoGetUploadURL(ctx, client, company)
					if err != nil {
						return err
					}

					// 2 Upload the file
					err = upload_file(uploadInfo.Url, contents, jsonFile.Name(), uploadInfo.Fields)
					if err != nil {
						return err
					}

					// 3 Notify using submit scan results
					blobInfo := openapi.BlobUploadInfo{
						Bucket: uploadInfo.Bucket,
						Key:    uploadInfo.Key,
					}

					results.BlobUpload = &blobInfo
				} else {
					var scans []openapi.SyndisInternalScanEventSyndisRiskScore
					err := json.Unmarshal([]byte(submitCmd_message), &scans)
					if err != nil {
						return err
					}
					events := openapi.BodySubmitScanResults_Events{}
					events.FromBodySubmitScanResultsEvents1(scans)
					results.Events = &events
				}

				resp, err := client.SubmitScanResults(ctx, scanName, results)

				if err != nil {
					return err
				}

				return openapi.CheckStatus(resp)
			default:
				return fmt.Errorf("unrecognised scan type %s", scanType)
			}

		},
	}
)

func validate_json_file(contents []byte) error {
	var scans []openapi.SyndisInternalScanEventSyndisRiskScore
	return json.Unmarshal(contents, &scans)
}
func createMultipartForm(contents []byte, filename string, fields map[string]string) (bytes.Buffer, *multipart.Writer, error) {
	var b bytes.Buffer

	w := multipart.NewWriter(&b)
	var fw io.Writer
	var err error

	for key, val := range fields {
		if fw, err = w.CreateFormField(key); err != nil {
			return b, nil, err
		}
		if _, err = io.Copy(fw, strings.NewReader(val)); err != nil {
			return b, nil, err
		}
	}

	if fw, err = w.CreateFormFile("file", filename); err != nil {
		return b, nil, err
	}
	if _, err = io.Copy(fw, bytes.NewBuffer(contents)); err != nil {
		return b, nil, err
	}

	w.Close()
	return b, w, nil
}

func upload_file(url string, contents []byte, filename string, fields map[string]string) error {
	byteBuffer, multiWriter, err := createMultipartForm(contents, filename, fields)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, &byteBuffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", multiWriter.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return errors.New(string(b))
	}
	return nil
}

func init() {

	submitCmd.Flags().StringVarP(&submitCmd_filename, "filename", "f", "", "JSON file to submit")
	submitCmd.Flags().StringVarP(&submitCmd_message, "message", "m", "", "JSON string to submit")

	// Want this, but no way to clear flag values between tests at the moment
	// https://github.com/spf13/cobra/issues/1180
	// submitCmd.MarkFlagsMutuallyExclusive("filename", "message")

	rootCmd.AddCommand(submitCmd)

}
