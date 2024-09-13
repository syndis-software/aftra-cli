/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	openapi "github.com/syndis-software/aftra-cli/pkg/openapi"
)

// getCmd represents the get command

var getConfigCmd = &cobra.Command{
	Use:   "config [scan-type] [scan-name]",
	Args:  cobra.MatchAll(cobra.RangeArgs(1, 2), cobra.OnlyValidArgs),
	Short: "Get the config for a scanner",
	Long: `Get the config for a scanner.

The output is suitable for being piped into a file for future use`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client := ctx.Value(clientKey).(*openapi.ClientWithResponses)
		scanType := args[0]
		scanName := ""

		if len(args) > 1 {
			scanName = args[1]
		}

		switch {
		case ScanType(scanType) == syndis:
			switch {
			case scanName == "":
				company := ctx.Value(companyKey).(string)
				configInfo, err := openapi.DoListCompanySyndisscans(ctx, company, client)
				if err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), configInfo)
				return nil
			default:
				configInfo, err := openapi.DoGetSyndisScanConfig(ctx, scanName, client)
				if err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), configInfo)
				return nil
			}

		default:
			return fmt.Errorf("unrecognised scan type %s", scanType)
		}

	},
}

func init() {
	getCmd.AddCommand(getConfigCmd)
}
