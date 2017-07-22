package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	infoCmd := &cobra.Command{
		Use:     "info",
		Aliases: []string{"status"},
		Short:   "Retrieve server status",
		Long:    "Retrieve server status",
		RunE:    info,
	}

	RootCmd.AddCommand(infoCmd)
}

func info(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("unimplemented")
}
