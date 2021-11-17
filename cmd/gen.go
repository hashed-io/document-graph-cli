package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/eoscanada/eos-go"
	"github.com/hashed-io/document-graph/docgraph"
	"github.com/spf13/cobra"
)

type FakeSocialUser struct {
	UserName  string `faker:"-"`
	DOB       string `faker:"date"`
	Timezone  string `faker:"timezone"`
	FirstName string `faker:"first_name"`
	Reference string `faker:"uuid_hyphenated"`
	Amount    string `faker:"amount_with_currency"`
	Joined    string `faker:"timestamp"`
}

type FakeSocialPost struct {
	ID      uint64
	Author  string
	Created string `faker:"timestamp"`
	Content string `faker:"paragraph"`
	Cost    eos.Asset
}

const charset = "abcdefghijklmnopqrstuvwxyz" + "12345"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandAccountName() string {
	return stringWithCharset(12, charset)
}

func constructContentItem(label, typeS string, val interface{}) docgraph.ContentItem {

	return docgraph.ContentItem{
		Label: label,
		Value: &docgraph.FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: docgraph.GetVariants().TypeID(typeS),
				Impl:   val}}}
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate graph data",
	Long:  "generate graph data",
	Run: func(cmd *cobra.Command, args []string) {

		for i := 0; i < 5; i++ {
			user := FakeSocialUser{
				UserName: RandAccountName(),
			}
			err := faker.FakeData(&user)
			if err != nil {
				panic(err)
			}

			systemGroup := make([]docgraph.ContentItem, 3)
			systemGroup[0] = constructContentItem("content_group_label", "string", "system")
			systemGroup[1] = constructContentItem("type", "name", "user")
			systemGroup[2] = constructContentItem("node_label", "string", "User: "+user.UserName)

			// username := constructContentItem("Username", "name", eos.Name(user.UserName))

			fmt.Printf("%+v\n", systemGroup)
		}

		// var ci docgraph.ContentItem
		// ci.Label = "label"
		// ci.Value = &docgraph.FlexValue{
		// 	BaseVariant: eos.BaseVariant{
		// 		TypeID: docgraph.FlexValueVariant.TypeID("name"),
		// 		Impl:   "nametype",
		// 	},
		// }

		// cg := make([]docgraph.ContentItem, 1)
		// cg[0] = ci
		// cgs := make([]docgraph.ContentGroup, 1)
		// cgs[0] = cg

		// doc := docgraph.Document{
		// 	Creator:       e.E().User,
		// 	ContentGroups: cgs,
		// }

		// docString, _ := json.MarshalIndent(doc, "", " ")
		// fmt.Println("doc: ", string(docString))

	},
}

func init() {
	genCmd.Flags().IntP("count", "n", 1, "number of documents to generate")

	RootCmd.AddCommand(genCmd)
}
