package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create documents or edges in the Document Graph",
}

func init() {
	RootCmd.AddCommand(createCmd)
}
