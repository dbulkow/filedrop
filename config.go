package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Write config file",
		Long:  "Write a config file with current values and defaults",
		RunE:  config,
	}

	RootCmd.AddCommand(configCmd)
}

func config(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("unimplemented")
}
