package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	storageCmd := &cobra.Command{
		Use:   "storage",
		Short: "Manage file storage",
		Long:  "Manage file storage",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List storage contents",
		Long:  "List detailed storage contents",
		RunE:  listStorage,
	}

	removeCmd := &cobra.Command{
		Use:     "remove <filename>...",
		Aliases: []string{"rm", "del", "delete"},
		Short:   "Remove file(s)",
		Long:    "Remove file(s)",
		RunE:    removeStorage,
	}

	storageCmd.AddCommand(listCmd)
	storageCmd.AddCommand(removeCmd)

	RootCmd.AddCommand(storageCmd)
}

func listStorage(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("unimplemented")
}

func removeStorage(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("unimplemented")
}

/*
.../storage/<hash of filename, time posted>/file
                                           /metadata
*/

type MetaData struct {
	Filename  string
	Hash      string
	Created   time.Time
	Expire    time.Time
	Owner     string
	MailAddrs []string
}

func (m *MetaData) Notify() {
	// iterate over MailAddrs
	// use smtp to provide html mail with link to server
}

func (m *MetaData) Expire() {
	if time.Now().After(m.Expire) {
		// remove hashdir
	}
}

func (m *MetaData) mkhash() {
}
