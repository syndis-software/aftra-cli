/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Aftra resources (eg tokens)",
	Long:  `Get Aftra resources via the API`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprint(cmd.ErrOrStderr(), "Error: must also specify a command")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
