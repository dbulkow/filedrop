package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List posted files",
		Long:    "List files posted to the server",
		RunE:    list,
	}

	RootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("unimplemented")
}
