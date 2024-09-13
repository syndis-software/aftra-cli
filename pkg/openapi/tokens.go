package openapi

import (
	"context"
)

func DoGetTokenInfo(ctx context.Context, client *ClientWithResponses) (*MaskedToken, error) {
	resp, err := client.GetTokenInfo(ctx)

	if err != nil {
		return nil, err
	}

	err = CheckStatus(resp)

	if err != nil {
		return nil, err
	}
	tokenInfo, err := ParseGetTokenInfoResponse(resp)

	if err != nil {
		return nil, err
	}
	return tokenInfo.JSON200, nil
}
