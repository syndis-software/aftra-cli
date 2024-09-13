package openapi

import (
	"context"
	"encoding/json"
	"fmt"
)

func DoGetSyndisScanConfig(ctx context.Context, configName string, client *ClientWithResponses) (string, error) {

	resp, err := client.GetSyndisConfigInfo(ctx, configName)

	if err != nil {
		return "", err
	}

	err = CheckStatus(resp)

	if err != nil {
		return "", err
	}

	config, err := ParseGetSyndisConfigInfoResponse(resp)

	if err != nil {
		return "", err
	}
	s, _ := json.Marshal(config.JSON200)
	return string(s), nil

}

func DoListCompanySyndisscans(ctx context.Context, companyPK string, client *ClientWithResponses) (string, error) {

	params := ListCompanySyndisscansParams{}

	resp, err := client.ListCompanySyndisscans(ctx, companyPK, &params)

	if err != nil {
		return "", err
	}

	err = CheckStatus(resp)

	if err != nil {
		return "", err
	}

	configs, err := ParseListCompanySyndisscansResponse(resp)
	if err != nil {
		return "", err
	}

	s := ""

	for _, entity := range configs.JSON200.Entities {
		x, _ := json.Marshal(entity)
		s += fmt.Sprintf("%s\n", string(x))
	}

	return s, nil
}
