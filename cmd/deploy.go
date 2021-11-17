package cmd

import (
	"fmt"
	"strconv"

	eostest "github.com/digital-scarcity/eos-go-test"
	"github.com/eoscanada/eos-go"
	"github.com/hashed-io/document-graph-cli/e"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy contracts",
	Long:  "deploy contracts",
	Run: func(cmd *cobra.Command, args []string) {
		// api := eos.New(viper.GetString("EosioEndpoint"))
		// ctx := context.Background()
		// contract := eos.AN(viper.GetString("DAOContract"))

		keyBag := &eos.KeyBag{}
		err := keyBag.ImportPrivateKey(e.E().X, eostest.DefaultKey())
		if err != nil {
			panic(fmt.Errorf("import private key failed: %v", err))
		}

		e.E().A.SetSigner(keyBag)

		contract, err := eostest.CreateAccountFromString(e.E().X, e.E().A, e.E().Contract.String(), eostest.DefaultKey())
		if err != nil {
			panic(fmt.Errorf("cannot create account: %v", err))
		}
		_, err = eostest.SetContract(e.Env.X, e.E().A, contract, "../document-graph/build/docs/docs.wasm", "../document-graph/build/docs/docs.abi")
		if err != nil {
			panic(fmt.Errorf("cannot set contract: %v", err))
		}
		for i := 1; i < 5; i++ {

			_, err := eostest.CreateAccountFromString(e.Env.X, e.E().A, "user"+strconv.Itoa(i), eostest.DefaultKey())
			if err != nil {
				panic(fmt.Errorf("cannot create account: %v", err))
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
