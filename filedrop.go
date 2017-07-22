package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "filedrop",
	Short: "File sharing service",
	Long:  "File sharing service",
}

func main() {
	var (
		config = os.Getenv("FILEDROP_CONFIG")
		server = os.Getenv("FILEDROP_SERVER")
	)

	RootCmd.PersistentFlags().StringVarP(&server, "server", "s", server, "Server address")
	RootCmd.PersistentFlags().StringVarP(&config, "config", "f", config, "Config file")

	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
