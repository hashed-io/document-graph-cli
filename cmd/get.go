package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get objects from the Document Graph",
}

func init() {
	RootCmd.AddCommand(getCmd)
}
