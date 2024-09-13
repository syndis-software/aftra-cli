/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	openapi "github.com/syndis-software/aftra-cli/pkg/openapi"
)

// tokenCmd represents the token command
var getTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Get token information for currently active token",
	Long: `Get token information for currently active token

Without args, this will output the full json text representing the token.

Supply "config" or "company" arguments to get escaped versions for 
use in setup scripts.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client := ctx.Value(clientKey).(*openapi.ClientWithResponses)

		tokenInfo, err := openapi.DoGetTokenInfo(ctx, client)

		if err != nil {
			return err
		}

		s, err := json.MarshalIndent(tokenInfo, "", "\t")

		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(s))

		return nil
	},
}

func init() {
	getCmd.AddCommand(getTokenCmd)
}
