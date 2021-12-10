package cmd

import (
	"context"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hashed-io/document-graph-cli/views"
	"github.com/hashed-io/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getEdgesCmd = &cobra.Command{
	Use:   "edges",
	Short: "print a list of all edges [ONLY use on very small graphs - no optimization]",
	Long:  "print a list of all edges [ONLY use on very small graphs - no optimization]",
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		edges, err := docgraph.GetAllEdges(ctx, api, eos.AN(viper.GetString("Contract")))
		if err != nil {
			panic(fmt.Errorf("cannot get all edges: %v", err))
		}

		edgesTable := views.EdgeTable(edges, false, false)
		edgesTable.SetStyle(simpletable.StyleCompactLite)
		fmt.Println("\n" + edgesTable.String() + "\n\n")
	},
}

func init() {
	getCmd.AddCommand(getEdgesCmd)
}
