package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/eoscanada/eos-go"
	"github.com/hashed-io/document-graph/docgraph"
	"github.com/spf13/viper"
)

func newDocumentTrx(ctx context.Context, api *eos.API,
	contract, creator eos.AccountName, data []byte) (docgraph.Document, error) {

	action := eos.ActN("create")

	var dump map[string]interface{}
	err := json.Unmarshal(data, &dump)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("unmarshal: %v", err)
	}

	dump["creator"] = creator

	actionBinary, err := api.ABIJSONToBin(ctx, contract, eos.Name(action), dump)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("api json to bin : %v", err)
	}

	actions := []*eos.Action{
		{
			Account: contract,
			Name:    action,
			Authorization: []eos.PermissionLevel{
				{Actor: creator, Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionDataFromHexData([]byte(actionBinary)),
		}}

	if viper.GetViper().GetBool("Vault") {
		pushEOSCActions(ctx, api, actions[0])
	} else {
		keyBag := &eos.KeyBag{}
		keys := viper.GetStringSlice("Keys")
		for _, key := range keys {
			keyBag.ImportPrivateKey(ctx, key)

		}
		api.SetSigner(keyBag)
		execWithRetry(ctx, api, actions)
	}

	lastDoc, err := docgraph.GetLastDocument(ctx, api, contract)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("get last document: %v", err)
	}
	return lastDoc, nil
}

func newDocumentFromFile(ctx context.Context, api *eos.API,
	contract, creator eos.AccountName, fileName string) (docgraph.Document, error) {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("readfile %v: %v", fileName, err)
	}

	return newDocumentTrx(ctx, api, contract, creator, data)
}

func newDocumentFromString(ctx context.Context, api *eos.API,
	contract, creator eos.AccountName, contentGroupsJson string) (docgraph.Document, error) {

	return newDocumentTrx(ctx, api, contract, creator, []byte(contentGroupsJson))
}

func execWithRetry(ctx context.Context, api *eos.API, actions []*eos.Action) (string, error) {
	trxId, err := exec(ctx, api, actions)

	if err != nil {
		if !strings.Contains(err.Error(), "deadline exceeded") {
			return string(""), err
		} else {
			attempts := 1
			for attempts < 3 {
				trxId, err = exec(ctx, api, actions)
				if err == nil {
					return trxId, nil
				}
				attempts++
			}
		}
		return string(""), err
	}
	return trxId, nil
}

func exec(ctx context.Context, api *eos.API, actions []*eos.Action) (string, error) {
	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(ctx, api); err != nil {
		return string(""), fmt.Errorf("error filling tx opts: %s", err)
	}

	tx := eos.NewTransaction(actions, txOpts)

	_, packedTx, err := api.SignTransaction(ctx, tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		return string(""), fmt.Errorf("error signing transaction: %s", err)
	}

	response, err := api.PushTransaction(ctx, packedTx)
	if err != nil {
		return string(""), fmt.Errorf("error pushing transaction: %s", err)
	}
	trxID := hex.EncodeToString(response.Processed.ID)
	return trxID, nil
}
