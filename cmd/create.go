/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>

*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// createCmd represents the create command

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create entities inside Aftra",
		Long: `Use the Aftra API to create things.

You'll need an API key to make this happen`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("must also specify a resource. eg opportunity")
		},
	}
)

func init() {
	rootCmd.AddCommand(createCmd)
}
