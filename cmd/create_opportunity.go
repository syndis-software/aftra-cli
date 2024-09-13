/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	openapi "github.com/syndis-software/aftra-cli/pkg/openapi"
)

// opportunityCmd represents the opportunity command
var (
	uid        string
	name       string
	score      int
	detailsStr string

	opportunityCmd = &cobra.Command{
		Use:          "opportunity",
		SilenceUsage: true,
		Short:        "Create internal opportunities inside Aftra",
		Long: `Use the Aftra API to create internal opportunities

	These will become part of the overall picture of your installation.
	You'll need an API key to make this happen`,
		RunE: func(cmd *cobra.Command, args []string) error {
			details, err := stringToMap(detailsStr)
			if err != nil {
				return err
			}

			opportunity := openapi.CreateOpportunity{
				Name:    name,
				Uid:     uid,
				Score:   openapi.OpportunityScore(score),
				Details: details,
			}

			ctx := cmd.Context()
			client := ctx.Value(clientKey).(*openapi.ClientWithResponses)
			company := ctx.Value(companyKey).(string)
			err = openapi.DoCreateOpportunity(ctx, client, company, opportunity)

			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s created\n", uid)
			return err
		},
	}
)

func stringToMap(str string) (map[string]openapi.CreateOpportunity_Details_AdditionalProperties, error) {
	result := make(map[string]openapi.CreateOpportunity_Details_AdditionalProperties)

	// split the string into key-value pairs
	pairs := strings.Split(str, ",")
	// loop through each key-value pair
	for _, pair := range pairs {
		// split the pair into key and value
		kv := strings.Split(pair, "=")

		// skip empty key-value pairs
		if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
			continue
		}

		// add the key-value pair to the map
		r := openapi.CreateOpportunity_Details_AdditionalProperties{}
		var err error
		if v_int, err := strconv.Atoi(kv[1]); err == nil {
			err = r.FromCreateOpportunityDetails1(v_int)
		} else {
			err = r.FromCreateOpportunityDetails0(kv[1])
		}
		if err != nil {
			return nil, err
		}
		result[kv[0]] = r
	}

	return result, nil
}

func init() {
	createCmd.AddCommand(opportunityCmd)
	opportunityCmd.Flags().StringVar(&uid, "uid", "", "Unique identifier for the opportunity")
	opportunityCmd.Flags().StringVar(&name, "name", "", "Name of the opportunity")
	opportunityCmd.Flags().IntVar(&score, "score", -1, "Risk score of the opportunity (critical (5), high (4), medium (3), low (2), info (1), none (0), unknown (-1))")
	opportunityCmd.Flags().StringVar(&detailsStr, "details", "", "Additional details. Comma separated key=value pairs.")
	opportunityCmd.MarkFlagRequired("uid")
	opportunityCmd.MarkFlagRequired("name")
}
