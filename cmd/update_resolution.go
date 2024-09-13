// Update Resolution usage
// aftra-cli update-resolution $UID $RESOLUTION --comment="" --due-date=2024-01-01

/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>
*/
package cmd

import (
	"fmt"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
	"github.com/syndis-software/aftra-cli/pkg/openapi"
	"golang.org/x/exp/slices"
)

func validateResolution(resolution string) (*openapi.OpportunityResolution, error) {

	possible := []openapi.OpportunityResolution{
		openapi.AcceptedRisk,
		openapi.Unacknowledged,
		openapi.FalsePositive,
		openapi.Resolved,
	}
	if !slices.Contains(possible, openapi.OpportunityResolution(resolution)) {
		return nil, fmt.Errorf("resolution must be one of %v. Got %s", possible, resolution)
	}
	final := openapi.OpportunityResolution(resolution)
	return &final, nil

}

func validateDueDateWithResolutionCheck(dueDate string, resolution openapi.OpportunityResolution) (*openapi_types.Date, error) {
	if resolution != openapi.AcceptedRisk {
		if dueDate != "" {
			return nil, fmt.Errorf("due date can only be set when setting resolution to accepted risk")
		}
		return nil, nil
	}
	t, err := time.Parse(openapi_types.DateFormat, dueDate)
	if err != nil {
		return nil, err
	}
	return &openapi_types.Date{t}, nil
}

var (
	commentStr string
	dueDateStr string
	dueDate    *openapi_types.Date

	updateResolutionsCmd = &cobra.Command{
		Use:   "resolution [uid] [status]",
		Short: "Update resolution",
		Long:  `Update a resolution of an opportunity.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			uidStr, resolutionStr := args[0], args[1]

			resolution, err := validateResolution(resolutionStr)
			if err != nil {
				return err
			}
			dueDate, err := validateDueDateWithResolutionCheck(dueDateStr, *resolution)
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			client := ctx.Value(clientKey).(*openapi.ClientWithResponses)
			company := ctx.Value(companyKey).(string)

			return openapi.DoPostResolution(
				ctx,
				client,
				company,
				uidStr,
				openapi.ResolutionUpdate{
					Comment:    &commentStr,
					DueDate:    dueDate,
					Resolution: openapi.OpportunityResolution(resolutionStr),
				},
			)

		},
	}
)

func init() {
	updateCmd.AddCommand(updateResolutionsCmd)
	updateResolutionsCmd.Flags().StringVar(&commentStr, "comment", "", "Comment associated with the resolution")
	updateResolutionsCmd.Flags().StringVar(&dueDateStr, "due-date", "", "Due date to set on the opportunity when changing resolution. Only valid when setting accepted risk.")

}
