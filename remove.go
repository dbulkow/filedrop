package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	removeCmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "del"},
		Short:   "Remove posted files",
		Long:    "Remove files posted to the server",
		RunE:    remove,
	}

	RootCmd.AddCommand(removeCmd)
}

func remove(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("unimplemented")
}
