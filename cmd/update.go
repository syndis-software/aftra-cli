/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update Aftra resources (eg resolutions)",
	Long:  `update Aftra resources via the API`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprint(cmd.ErrOrStderr(), "Error: must also specify a command")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
