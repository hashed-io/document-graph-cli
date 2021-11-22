package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/eoscanada/eos-go"
	"github.com/hashed-io/document-graph-cli/e"
	"github.com/hashed-io/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FakeSocialUser struct {
	ID        uint64 `faker:"-"`
	UserName  string `faker:"-"`
	DOB       string `faker:"date"`
	Timezone  string `faker:"timezone"`
	FirstName string `faker:"first_name"`
	Reference string `faker:"uuid_hyphenated"`
	Amount    string `faker:"amount_with_currency"`
	Joined    string `faker:"timestamp"`
}

type FakeSocialPost struct {
	ID      uint64    `faker:"-"`
	Author  string    `faker:"-"`
	Created string    `faker:"timestamp"`
	Title   string    `faker:"sentence"`
	Content string    `faker:"paragraph"`
	Cost    eos.Asset `faker:"-"`
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

func randInt(min int, max int) int {
	return min + seededRand.Intn(max-min)
}

func constructContentItem(label, typeS string, val interface{}) docgraph.ContentItem {

	return docgraph.ContentItem{
		Label: label,
		Value: &docgraph.FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: docgraph.GetVariants().TypeID(typeS),
				Impl:   val}}}
}

var userCount, postCount int

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate graph data",
	Long:  "generate graph data",
	Run: func(cmd *cobra.Command, args []string) {

		os.MkdirAll(viper.GetString("GeneratedDir"), os.ModePerm)

		users := make([]FakeSocialUser, userCount)
		posts := make([]FakeSocialPost, postCount)

		for i := 0; i < len(users); i++ {
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

			detailsGroup := make([]docgraph.ContentItem, 8)
			detailsGroup[0] = constructContentItem("content_group_label", "string", "details")
			detailsGroup[1] = constructContentItem("username", "name", user.UserName)
			detailsGroup[2] = constructContentItem("dob", "string", user.DOB)
			detailsGroup[3] = constructContentItem("timezone", "string", user.Timezone)
			detailsGroup[4] = constructContentItem("firstname", "string", user.FirstName)
			detailsGroup[5] = constructContentItem("reference", "string", user.Reference)
			detailsGroup[6] = constructContentItem("stake", "string", user.Amount)
			detailsGroup[7] = constructContentItem("joined", "string", user.Joined)

			cgs := make([]docgraph.ContentGroup, 2)
			cgs[0] = systemGroup
			cgs[1] = detailsGroup

			doc := docgraph.Document{}
			doc.ContentGroups = cgs

			docString, _ := json.MarshalIndent(doc, "", " ")

			filename := viper.GetString("GeneratedDir") + "/user_" + user.UserName + ".json"
			_ = ioutil.WriteFile(filename, docString, 0644)

			contentGroupsString, _ := json.Marshal(doc.ContentGroups)
			csgString := "{\"content_groups\":" + string(contentGroupsString) + "}"
			// fmt.Printf("%+v\n\n", csgString)

			createdDoc, err := newDocumentFromString(e.E().X, e.E().A, e.E().Contract, e.E().User, csgString)
			if err != nil {
				panic(err)
			}

			fmt.Println("Successfully created document.	ID 	: " + strconv.Itoa(int(createdDoc.ID)))
			user.ID = createdDoc.ID
			users[i] = user
		}

		for i := 0; i < len(posts); i++ {
			post := FakeSocialPost{}

			err := faker.FakeData(&post)
			if err != nil {
				panic(err)
			}

			author := users[randInt(0, len(users))]
			systemGroup := make([]docgraph.ContentItem, 3)
			systemGroup[0] = constructContentItem("content_group_label", "string", "system")
			systemGroup[1] = constructContentItem("type", "name", "post")
			systemGroup[2] = constructContentItem("node_label", "string", "Post: "+post.Title)

			detailsGroup := make([]docgraph.ContentItem, 5)
			detailsGroup[0] = constructContentItem("content_group_label", "string", "details")
			detailsGroup[1] = constructContentItem("author", "name", author.UserName) //users[rand.Intn(len(users))].UserName)
			detailsGroup[2] = constructContentItem("title", "string", post.Title)
			detailsGroup[3] = constructContentItem("content", "string", post.Content)

			assetAmount := strconv.Itoa(randInt(0, 1000)) + ".0000 USD"
			cost, _ := eos.NewAssetFromString(assetAmount)

			detailsGroup[4] = constructContentItem("cost", "asset", cost)

			cgs := make([]docgraph.ContentGroup, 2)
			cgs[0] = systemGroup
			cgs[1] = detailsGroup

			doc := docgraph.Document{}
			doc.ContentGroups = cgs

			docString, _ := json.MarshalIndent(doc, "", " ")
			// fmt.Printf("%+v\n\n", string(docString))
			_ = ioutil.WriteFile(viper.GetString("GeneratedDir")+"/post_"+strings.ReplaceAll(post.Title, " ", "_")+"json", docString, 0644)

			contentGroupsString, _ := json.Marshal(doc.ContentGroups)
			csgString := "{\"content_groups\":" + string(contentGroupsString) + "}"
			//fmt.Printf("%+v\n\n", csgString)

			createdDoc, err := newDocumentFromString(e.E().X, e.E().A, e.E().Contract, e.E().User, csgString)
			if err != nil {
				panic(err)
			}

			fmt.Println("Successfully created document.	ID 	: " + strconv.Itoa(int(createdDoc.ID)))

			post.ID = createdDoc.ID

			_, err = docgraph.CreateEdge(e.E().X, e.E().A, e.E().Contract, e.E().User, author.ID, post.ID, eos.Name("authored"))
			if err != nil {
				panic(err)
			}
			fmt.Println("Successfully created edge.")

			for i := 0; i < 4; i++ {
				_, err = docgraph.CreateEdge(e.E().X, e.E().A, e.E().Contract, e.E().User, users[randInt(0, len(users))].ID, post.ID, eos.Name("liked"))
				if err != nil {
					fmt.Println("liked failed")
				} else {
					fmt.Println("post is liked")
				}
			}

			posts[i] = post
		}
	},
}

func init() {
	genCmd.Flags().IntVarP(&userCount, "users", "u", 1, "number of documents to generate")
	genCmd.Flags().IntVarP(&postCount, "posts", "p", 1, "number of posts to generate")

	RootCmd.AddCommand(genCmd)
}
