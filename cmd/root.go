/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	openapi "github.com/syndis-software/aftra-cli/pkg/openapi"
)

// rootCmd represents the base command when called without any subcommands
var (
	// Used for flags
	rootCmd = &cobra.Command{
		Use:          "aftra-cli",
		SilenceUsage: true,
		Short:        "CLI for the Aftra API",
		Long:         `CLI for using the AFTRA API`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(cmd.ErrOrStderr(), "Error: must also specify a command")
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			api_key := viper.GetString("api_token")
			company := viper.GetString("company")
			host := viper.GetString("host")

			ctx := cmd.Context()
			doer := ctx.Value(doerKey).(openapi.HttpRequestDoer)

			apiKeyIntercept, _ := openapi.NewSecurityProviderApiKey("x-api-key", api_key)
			client, _ := openapi.NewClientWithResponses(host, openapi.WithRequestEditorFn(apiKeyIntercept.Intercept), openapi.WithHTTPClient(doer))
			ctx = context.WithValue(ctx, clientKey, client)
			ctx = context.WithValue(ctx, companyKey, company)
			cmd.SetContext(ctx)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context, doer openapi.HttpRequestDoer) {
	ctx = context.WithValue(ctx, doerKey, doer)
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().String("host", "https://app.aftra.io", "Aftra host (AFTRA_HOST)")
	rootCmd.PersistentFlags().String("company", "", "Company ID. Should look like Company-XXXX (AFTRA_COMPANY)")
	viper.BindPFlag("company", rootCmd.PersistentFlags().Lookup("company"))
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))

	viper.SetEnvPrefix("aftra")
	viper.AutomaticEnv()

}
