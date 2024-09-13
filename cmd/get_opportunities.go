// Get opportunities usage
// aftra-cli get opportunities --updated-since=DT --limit=100

/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
	openapi "github.com/syndis-software/aftra-cli/pkg/openapi"
)

// getCmd represents the get command

func getTime(s string) (time.Time, error) {
	// Taken from proposed implementation of parsing times in cobra
	// https://github.com/spf13/pflag/pull/348/files
	s = strings.TrimSpace(s)
	formats := []string{time.RFC3339Nano, time.RFC1123Z}
	for _, f := range formats {
		v, err := time.Parse(f, s)
		if err != nil {
			continue
		}
		return v, nil
	}

	formatsString := ""
	for i, f := range formats {
		if i > 0 {
			formatsString += ", "
		}
		formatsString += fmt.Sprintf("`%s`", f)
	}

	return time.Time{}, fmt.Errorf("invalid time format `%s` must be one of: %s", s, formatsString)
}

func validateLimit(l int) error {
	// limit can either be -1 or 0 < l <= 1000
	if (limit == -1) || (limit > 0 && limit <= 1000) {
		return nil
	}
	return fmt.Errorf("limit should be -1 (everything) or less than 1000: %d", l)
}

var (
	limit        int
	updatedSince string

	getOpportunitiesCmd = &cobra.Command{
		Use:   "opportunities ",
		Short: "Get opportunities",
		Long: `Get filtered opportunities.

Output is JSON format`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			client := ctx.Value(clientKey).(*openapi.ClientWithResponses)
			company := ctx.Value(companyKey).(string)

			lastUpdatedGte, err := getTime(updatedSince)

			if err != nil {
				return err
			}
			err = validateLimit(limit)

			if err != nil {
				return err
			}

			var order openapi.SearchOpportunitiesApiCompaniesCompanyPkOpportunitiesV3GetParamsOrder = "asc"
			var sort openapi.SortOptions = "timestamp_last_updated"

			var startFrom = 0
			var totalFetched = 0

			//  if limit is -1 (unset), we want all opportunities
			//  if limit is <1000 we want one page of opportunities up to that count
			var batchSize int
			var totalForSearch = -1
			if limit == -1 {
				batchSize = 1000
			} else {
				batchSize = limit
			}

			for totalForSearch == -1 || totalFetched < totalForSearch {
				params := openapi.SearchOpportunitiesApiCompaniesCompanyPkOpportunitiesV3GetParams{
					TimestampLastUpdatedGte: &openapi_types.Date{lastUpdatedGte},
					Sort:                    &sort,
					Order:                   &order,
					Limit:                   &batchSize,
					StartFrom:               &startFrom,
				}

				opportunities, err := openapi.DoSearchOpportunities(ctx, client, company, params)

				if err != nil {
					return err
				}

				for _, oppo := range opportunities.Opportunities {
					txt, err := json.Marshal(oppo)
					if err != nil {
						return err
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", txt)
					totalFetched += 1
					startFrom += 1
				}
				if totalForSearch == -1 {
					if limit == -1 {
						totalForSearch = opportunities.Total
					} else {
						totalForSearch = int(math.Min(float64(limit), float64(opportunities.Total)))
					}

				}
			}

			return nil
		},
	}
)

func init() {
	getCmd.AddCommand(getOpportunitiesCmd)
	getOpportunitiesCmd.Flags().IntVar(&limit, "limit", -1, "Max number of opportunities to retrieve")
	getOpportunitiesCmd.Flags().StringVar(&updatedSince, "updated-since", "2020-01-01T00:00:00Z", "Only fetch opportunities updated since")

}
