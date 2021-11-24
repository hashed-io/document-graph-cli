package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/hashed-io/document-graph-cli/util"
	"github.com/hashed-io/document-graph/docgraph"
	"github.com/manifoldco/promptui"
	"github.com/patrickmn/go-cache"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
	"go.uber.org/zap"
)

func cleanString(input string) string {
	input = strings.Replace(input, "\n", "", -1)

	if len(input) > 65 {
		return input[:40]
	}
	return input
}

func printContentGroups(p *Page) {

	fmt.Println("ContentGroups")
	for _, contentGroup := range p.Primary.ContentGroups {
		fmt.Println("  ContentGroup")

		for _, content := range contentGroup {
			fmt.Print("    ")
			fmt.Printf("%-35v", cleanString(content.Label))
			fmt.Printf("%-65v\n", cleanString(content.Value.String()))
		}
	}
	fmt.Println()
}

// func printEdges(ctx context.Context, api *eos.API, p *Page) {

// 	colorCyan := "\033[36m"
// 	colorReset := "\033[0m"
// 	colorRed := "\033[31m"

// 	if len(p.ToEdges) > 0 {

// 		sort.SliceStable(p.ToEdges, func(i, j int) bool {
// 			return p.ToEdges[i].CreatedDate.Before(p.ToEdges[j].CreatedDate.Time)
// 		})

// 		toEdgesTable := views.EdgeTable(p.ToEdges, true, false)
// 		toEdgesTable.SetStyle(simpletable.StyleCompactLite)

// 		fmt.Print(string(colorCyan), "\n                                                                                           to this node  ---------->")
// 		fmt.Println(string(colorReset))
// 		fmt.Println("\n" + toEdgesTable.String() + "\n\n")
// 	}

// 	if len(p.FromEdges) > 0 {

// 		sort.SliceStable(p.FromEdges, func(i, j int) bool {
// 			return p.FromEdges[i].CreatedDate.Before(p.FromEdges[j].CreatedDate.Time)
// 		})

// 		fromEdgesTable := views.EdgeTable(p.FromEdges, true, true)
// 		fromEdgesTable.SetStyle(simpletable.StyleCompactLite)

// 		fmt.Print(string(colorCyan), "\n                                    from this node  ---------->")
// 		fmt.Println(string(colorReset))
// 		fmt.Println("\n" + fromEdgesTable.String() + "\n\n")
// 		// fmt.Print(string(colorCyan), "\n  <---   use 'rev <edge-name>' for reverse                                        use 'fwd <edge-name>' for forward --->")
// 		// fmt.Print(string(colorCyan), "\n  <edge-name> defaults to first in list if left blank")
// 		fmt.Println(string(colorReset))

// 	} else {
// 		fmt.Print(string(colorRed), "\n                                                                                                        ----------")
// 		fmt.Print(string(colorRed), "\n                                                                                                        |  END   |")
// 		fmt.Print(string(colorRed), "\n                                                                                                        ----------")
// 		fmt.Println(string(colorReset))
// 	}
// }

func printDocument(ctx context.Context, api *eos.API, p *Page) {
	fmt.Println("Document Details")

	fmt.Println()
	output := []string{
		fmt.Sprintf("ID|%v", strconv.Itoa(int(p.Primary.ID))),
		fmt.Sprintf("Hash|%v", p.Primary.Hash.String()),
		fmt.Sprintf("Creator|%v", string(p.Primary.Creator)),
		fmt.Sprintf("Created Date|%v", p.Primary.CreatedDate.Time.Format("2006 Jan 02 15:04:05")),
	}

	fmt.Println(columnize.SimpleFormat(output))
	fmt.Println()
	printContentGroups(p)
	// printEdges(ctx, api, p)
}

// Page ...
type Page struct {
	Primary     docgraph.Document
	FromEdges   []docgraph.Edge
	ToEdges     []docgraph.Edge
	EdgePrompts []edgeChoice
	Graph       image.Image
}

// func (p *Page) getKey() string {
// 	return p.Primary.Hash.String()
// }

func getLabel(d *docgraph.Document) string {
	documentLabel, _ := d.GetContent("node_label")
	if documentLabel != nil {
		if len(documentLabel.String()) <= 40 {
			return documentLabel.String()
		}
		return documentLabel.String()[:40]
	}
	return "unlabeled"
}

func getType(d *docgraph.Document) string {
	documentType, _ := d.GetContent("type")
	if documentType != nil {
		return documentType.String()
	}
	return "untyped"
}

type edgeChoice struct {
	Name        eos.Name
	Forward     bool
	To          string
	ToType      string
	ToLabel     string
	CreatedDate eos.BlockTimestamp
	Column1     string // First 69 characters
	Column2     string //
}

func (e *edgeChoice) GetWidth() int {
	return 25
}

func getPage(ctx context.Context, api *eos.API, pageCache, documentCache *cache.Cache, contract eos.AccountName, id uint64) (Page, error) {
	pager, found := pageCache.Get(strconv.Itoa(int(id)))
	if found {
		return pager.(Page), nil
	}

	var err error
	page := Page{}
	page.Primary, err = getDocument(ctx, api, documentCache, contract, id)
	if err != nil {
		zlog.Error("document not found", zap.Uint64("document-id", id))
		return page, err
	}

	page.FromEdges, err = docgraph.GetEdgesFromDocument(ctx, api, contract, page.Primary)
	if err != nil {
		log.Println("ERROR: Cannot get edges from document: ", err)
	}

	page.ToEdges, err = docgraph.GetEdgesToDocument(ctx, api, contract, page.Primary)
	if err != nil {
		log.Println("ERROR: Cannot get edges to document: ", err)
	}

	page.EdgePrompts = make([]edgeChoice, len(page.FromEdges)+len(page.ToEdges))

	for i, edge := range page.FromEdges {

		document, err := getDocument(ctx, api, documentCache, contract, edge.ToNode)
		if err != nil {
			zlog.Error("document not found", zap.Uint64("document-id", id))
			return Page{}, err
		}

		page.EdgePrompts[i] = edgeChoice{
			Name:        edge.EdgeName,
			Forward:     true,
			To:          strconv.Itoa(int(edge.ToNode)),
			ToType:      getType(&document),
			ToLabel:     getLabel(&document),
			CreatedDate: edge.CreatedDate,
			Column1:     string("                                                                          [x]"),
			Column2: string(" ---> " + fmt.Sprintf("%-12v", string(edge.EdgeName)) +
				" ---> " + getLabel(&document) + " (" + getType(&document) + ")"),
		}
	}

	for i, edge := range page.ToEdges {

		document, err := getDocument(ctx, api, documentCache, contract, edge.FromNode)
		if err != nil {
			zlog.Error("document not found", zap.Uint64("document-id", id))
			return Page{}, err
		}

		page.EdgePrompts[i+len(page.FromEdges)] = edgeChoice{
			Name:        edge.EdgeName,
			Forward:     false,
			To:          strconv.Itoa(int(edge.FromNode)),
			ToType:      getType(&document),
			ToLabel:     getLabel(&document),
			CreatedDate: edge.CreatedDate,
			Column1: string(
				fmt.Sprintf("%49v", getLabel(&document)+
					" ("+
					getType(&document)) +
					") ---> " +
					fmt.Sprintf("%12v", string(edge.EdgeName)) +
					" ---> [x]"),
			Column2: string("                                                                          "),
		}
	}

	sort.SliceStable(page.EdgePrompts, func(i, j int) bool {
		return page.EdgePrompts[i].CreatedDate.Before(page.EdgePrompts[j].CreatedDate.Time)
	})

	pageCache.Set(page.Primary.Hash.String(), page, cache.DefaultExpiration)

	return page, nil
}

func getDocument(ctx context.Context, api *eos.API, c *cache.Cache, contract eos.AccountName, id uint64) (docgraph.Document, error) {

	documenter, found := c.Get(strconv.Itoa(int(id)))
	if found {
		return documenter.(docgraph.Document), nil
	}

	document, err := docgraph.LoadDocumentByID(ctx, api, contract, id)
	if err != nil {
		zlog.Error("document not found", zap.Uint64("document-id", id))
		return docgraph.Document{}, err
	}
	c.Set(document.Hash.String(), document, cache.DefaultExpiration)

	return document, nil
}

func loadCache(ctx context.Context, api *eos.API, pages, documents *cache.Cache, contract eos.AccountName, startingNode uint64) error {

	go func() error {
		page, err := getPage(ctx, api, pages, documents, contract, startingNode)
		if err != nil {
			zlog.Error("error getting page", zap.Uint64("document-id", startingNode), zap.Error(err))
			return fmt.Errorf("error getting page: %v", err)
		}

		for _, edge := range page.ToEdges {
			_, err := getPage(ctx, api, pages, documents, contract, edge.ToNode)
			if err != nil {
				zlog.Error("error getting page", zap.Uint64("document-id", edge.ToNode), zap.Error(err))
				return fmt.Errorf("error getting page: %v", err)
			}
		}

		for _, edge := range page.FromEdges {
			_, err := getPage(ctx, api, pages, documents, contract, edge.ToNode)
			if err != nil {
				zlog.Error("error getting page", zap.Uint64("document-id", edge.ToNode), zap.Error(err))
				return fmt.Errorf("error getting page: %v", err)
			}
		}
		return nil
	}()
	return nil
}

var getDocumentCmd = &cobra.Command{
	Use:   "document [hash | shorty]",
	Short: "retrieve document details and navigate the graph",
	Long:  "retrieve the detailed content within a document",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		contract := eos.AN(viper.GetString("Contract"))

		var id uint64

		// if last==true OR no argument, use the last document
		if viper.GetBool("get-document-cmd-last") || len(args) == 0 {
			lastDocument, err := docgraph.GetLastDocument(ctx, api, contract)
			if err != nil {
				return fmt.Errorf("cannot get last document: %v", err)
			}
			id = lastDocument.ID
		}

		// if getting a document with JSON, just print it out and exit
		if viper.GetBool("get-document-cmd-json") {
			document, err := util.Get(ctx, api, contract, strconv.Itoa(int(id)))
			if err != nil {
				return fmt.Errorf("cannot find document with hash: %v %v", id, err)
			}

			docJson, err := json.Marshal(document)
			if err != nil {
				return fmt.Errorf("cannot marshall document to JSON: %v %v", args[0], err)
			}

			fmt.Println(string(pretty.Color(pretty.Pretty(docJson), nil)))
			return nil
		}

		if len(args) == 1 {
			idInt, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("parameter must be formatted as int: %v %v", args[0], err)
			}

			id = uint64(idInt)
		}

		pages := cache.New(5*time.Minute, 10*time.Minute)
		documents := cache.New(5*time.Minute, 10*time.Minute)

		for {

			page, err := getPage(ctx, api, pages, documents, contract, id)
			if err != nil {
				zlog.Error("error getting page", zap.Uint64("document-id", id), zap.Error(err))
				return fmt.Errorf("error getting page: %v", err)
			}

			loadCache(ctx, api, pages, documents, contract, id)

			printDocument(ctx, api, &page)

			if viper.GetBool("get-document-cmd-navigate") {
				fmt.Println("                          ")
				templates := &promptui.SelectTemplates{
					Label:    "{{ . }}  {{ \"from: node_label (type) --->    edge_name ---> [X] ---> edge_name    ---> to: node_label (type) \" | faint }}",
					Active:   "{{ .Column1 | cyan}}{{ .Column2 | cyan}}",
					Inactive: "{{ .Column1 }}{{ .Column2 }}",
					Selected: "{{ \"  -- selected -- \" | yellow}}{{ .Column1 | yellow}}{{ .Column2 | yellow}} {{ \"  -- selected -- \" | yellow}}",
				}

				searcher := func(input string, index int) bool {
					pepper := page.EdgePrompts[index]
					name := strings.Replace(strings.ToLower(string(pepper.Name)), " ", "", -1)
					input = strings.Replace(strings.ToLower(input), " ", "", -1)

					return strings.Contains(name, input)
				}

				prompt2 := promptui.Select{
					Label:     "Select an edge to navigate:",
					Items:     page.EdgePrompts,
					Templates: templates,
					Size:      len(page.FromEdges) + len(page.ToEdges),
					Searcher:  searcher,
				}

				i, _, err := prompt2.Run()
				if err != nil {
					return fmt.Errorf("prompt failed %v", err)
				}

				fmt.Println()
				fmt.Println("-------------------------------------------------------------------------------------------------")
				fmt.Println()

				interimId, _ := strconv.Atoi(page.EdgePrompts[i].To)
				id = uint64(interimId)
			} else {
				return nil
			}
		}
	},
}

func init() {
	getDocumentCmd.Flags().BoolP("last", "l", false, "retrieve the most recently created document")
	getDocumentCmd.Flags().BoolP("json", "j", false, "print the document to the terminal in JSON and exit")
	getDocumentCmd.Flags().BoolP("navigate", "n", true, "show document edges and allow interactive graph navigation")
	getCmd.AddCommand(getDocumentCmd)
}
