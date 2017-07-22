package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	storageCmd := &cobra.Command{
		Use:   "storage",
		Short: "Dump storage contents",
		Long:  "Dump detailed storage contents",
		RunE:  storage,
	}

	RootCmd.AddCommand(storageCmd)
}

func storage(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("unimplemented")
}
