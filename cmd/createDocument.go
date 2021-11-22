package cmd

import (
	"fmt"

	"github.com/hashed-io/document-graph-cli/e"
	"github.com/hashed-io/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createDocumentCmd = &cobra.Command{
	Use:   "document <content-groups>",
	Short: "create a new document in the Document Graph",
	// Args:  cobra.RangeArgs(1, 1),
	RunE: func(cmd *cobra.Command, args []string) error {

		var createdDoc docgraph.Document
		var err error
		if len(viper.GetString("create-document-cmd-content")) > 0 {
			// zlog.Fatal("User passed content on terminal -- NOT IMPLEMENTED YET")
			createdDoc, err = newDocumentFromString(e.E().X, e.E().A, e.E().Contract, e.E().User, viper.GetString("create-document-cmd-content"))
			if err != nil {
				return fmt.Errorf("cannot create document from string: %v", err)
			}
		}

		if len(viper.GetString("create-document-cmd-file")) == 0 {
			zlog.Fatal("Either --content or --file parameter is required")
		}

		createdDoc, err = newDocumentFromFile(e.E().X, e.E().A, e.E().Contract, e.E().User, viper.GetString("create-document-cmd-file"))
		if err != nil {
			return fmt.Errorf("cannot create document from file: %v %v", viper.GetString("create-document-cmd-file"), err)
		}

		fmt.Println("Successfully created document: ", createdDoc.ID)
		return nil
	},
}

func init() {
	createDocumentCmd.Flags().StringP("content", "", "", "content groups (json)")
	createDocumentCmd.Flags().StringP("file", "f", "", "file with content groups")
	createCmd.AddCommand(createDocumentCmd)
}
