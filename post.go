package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	postCmd := &cobra.Command{
		Use:   "post",
		Short: "Post a file",
		Long:  "Post a file to the server",
		RunE:  post,
	}

	RootCmd.AddCommand(postCmd)
}

func post(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("unimplemented")
}
