package openapi

import (
	"context"
)

func DoGetUploadURL(ctx context.Context, client *ClientWithResponses, companyPk string) (*BlobSignedUploadURLResponse, error) {

	params := GetUploadURLParams{}
	resp, err := client.GetUploadURLWithResponse(ctx, companyPk, &params)

	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil

}
