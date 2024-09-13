package openapi

import (
	"context"
)

func DoCreateOpportunity(ctx context.Context, client *ClientWithResponses, companyPk string, opportunity CreateOpportunity) error {

	params := CreateOpportunityApiCompaniesCompanyPkOpportunitiesPostParams{}
	resp, err := client.CreateOpportunityApiCompaniesCompanyPkOpportunitiesPost(
		ctx,
		companyPk,
		&params,
		CreateOpportunityApiCompaniesCompanyPkOpportunitiesPostJSONRequestBody(opportunity),
	)
	if err != nil {
		return err
	}

	return CheckStatus(resp)
}

func DoCreateExternalOpportunity(ctx context.Context, client *ClientWithResponses, companyPk string, opportunity CreateExternalOpportunity) error {

	params := CreateExternalOpportunityApiCompaniesCompanyPkOpportunitiesExternalPostParams{}
	resp, err := client.CreateExternalOpportunityApiCompaniesCompanyPkOpportunitiesExternalPost(
		ctx,
		companyPk,
		&params,
		CreateExternalOpportunityApiCompaniesCompanyPkOpportunitiesExternalPostJSONRequestBody(opportunity),
	)
	if err != nil {
		return err
	}

	return CheckStatus(resp)
}

func DoSearchOpportunities(ctx context.Context, client *ClientWithResponses, companyPk string, params SearchOpportunitiesApiCompaniesCompanyPkOpportunitiesV3GetParams) (*SearchedOpportunitiesResponse, error) {

	resp, err := client.SearchOpportunitiesApiCompaniesCompanyPkOpportunitiesV3Get(ctx, companyPk, &params)
	if err != nil {
		return nil, err
	}

	err = CheckStatus(resp)

	if err != nil {
		return nil, err
	}
	opportunities, err := ParseSearchOpportunitiesApiCompaniesCompanyPkOpportunitiesV3GetResponse(resp)

	if err != nil {
		return nil, err
	}
	return opportunities.JSON200, nil

}

func DoPostResolution(ctx context.Context, client *ClientWithResponses, companyPk string, opportunityUid string, update ResolutionUpdate) error {

	params := PostUpdateOpportunityResolutionParams{}
	resp, err := client.PostUpdateOpportunityResolution(ctx, companyPk, opportunityUid, &params, PostUpdateOpportunityResolutionJSONRequestBody(update))

	if err != nil {
		return err
	}

	err = CheckStatus(resp)
	if err != nil {
		return err
	}
	return nil
}
