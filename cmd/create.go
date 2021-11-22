package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create objects in the Document Graph",
}

func init() {
	RootCmd.AddCommand(createCmd)
}
