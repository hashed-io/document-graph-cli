package cmd

import (
	"fmt"
	"strconv"

	"github.com/hashed-io/document-graph-cli/e"
	"github.com/hashed-io/document-graph/docgraph"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var content, file string

var createDocumentCmd = &cobra.Command{
	Use:   "document <content-groups>",
	Short: "create a new document in the Document Graph",
	RunE: func(cmd *cobra.Command, args []string) error {

		var createdDoc docgraph.Document
		var err error
		if len(content) > 0 {
			createdDoc, err = newDocumentFromString(e.E().X, e.E().A, e.E().Contract, e.E().User, content)
			if err != nil {
				zlog.Error("cannot create document from string", zap.String("content", content), zap.Error(err))
				return fmt.Errorf("cannot create document from string: %v", err)
			}
		}

		if len(file) == 0 {
			zlog.Fatal("Either --content or --file parameter is required")
		}

		createdDoc, err = newDocumentFromFile(e.E().X, e.E().A, e.E().Contract, e.E().User, file)
		if err != nil {
			zlog.Error("cannot create document from file", zap.String("content", content), zap.String("file", file))
			return fmt.Errorf("cannot create document from file: %v %v %v", file, content, err)
		}

		zlog.Debug("created document", zap.String("document-id", strconv.Itoa(int(createdDoc.ID))))
		return nil
	},
}

func init() {
	createDocumentCmd.Flags().StringVarP(&content, "content", "", "", "content groups (json)")
	createDocumentCmd.Flags().StringVarP(&file, "file", "f", "", "file with content groups")
	createCmd.AddCommand(createDocumentCmd)
}
