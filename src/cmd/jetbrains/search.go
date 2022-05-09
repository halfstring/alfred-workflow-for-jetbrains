package jetbrains

import (
	"bufio"
	"encoding/json"
	"fmt"
	"jetbrains-workflow/golang/tools"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Result struct {
	Items     []Item    `json:"items"`
	Variables Variables `json:"variables"`
}

type Item struct {
	BasicItem
	Arg       string    `json:"arg"`
	Valid     bool      `json:"valid"`
	ItemText  ItemText  `json:"text"`
	Type      string    `json:"type"`
	Variables Variables `json:"variables"`
}

type ItemText struct {
	Copy string `json:"copy"`
}

type Variables struct {
	JbBin string `json:"jb_bin"`
}

type EmptyResult struct {
	BasicItems []BasicItem `json:"items"`
}

type BasicItem struct {
	Title        string `json:"title"`
	SubTitle     string `json:"subtitle"`
	Match        string `json:"match"`
	AutoComplete string `json:"autocomplete"`
}

func init() {
	//cmd.RootCmd.AddCommand(SearchCommand)

	SearchCommand.Flags().String("keyword", "", "Please input your keyword to continue.")
	SearchCommand.Flags().String("plate", "", "严格匹配，默认不开启.")
}

var SearchCommand = &cobra.Command{
	Use:   "project/search",
	Short: "filter project",
	Run: func(cmd *cobra.Command, args []string) {
		keyword, _ := cmd.Flags().GetString("keyword")
		plate, _ := cmd.Flags().GetString("plate")
		plate = strings.TrimSpace(plate)

		keyword = strings.TrimSpace(keyword)

		bin, _ := cmd.Root().Flags().GetString("bin")
		bin = strings.TrimSpace(bin)

		dataDir, _ := tools.GetDataDir(cmd)
		filePath := dataDir + "/projects.list"

		_, b := tools.IsFile(filePath)
		if !b {
			f, _ := os.Create(filePath)
			defer f.Close()
		}

		f, err := os.OpenFile(filePath, os.O_APPEND, 0666)
		if err != nil {
			defer f.Close()

			fmt.Println("open file failed.", err)
			os.Exit(1)
		}

		fileRead, _ := os.Open(filePath)
		defer fileRead.Close()

		buf := bufio.NewScanner(fileRead)
		ItemsList := make([]Item, 0)
		for {
			if !buf.Scan() {
				break
			}
			line := buf.Text()

			p := strings.Split(line, "\t")
			if len(p) != 2 && len(p) != 3 {
				continue
			}

			if keyword == "" || strings.Contains(strings.ToLower(p[0]), strings.ToLower(keyword)) {
				if len(p) == 3 && plate != "" {
					if strings.ToLower(p[2]) != strings.ToLower(plate) {
						continue
					}
				}

				ItemsList = append(ItemsList, Item{
					BasicItem: BasicItem{
						Title:        p[0],
						SubTitle:     p[1],
						Match:        p[0],
						AutoComplete: p[0],
					},
					Arg:   p[1],
					Valid: true,
					Type:  "default",
					ItemText: ItemText{
						Copy: "",
					},
					Variables: Variables{
						JbBin: bin,
					},
				})

			}
		}

		if len(ItemsList) < 1 {
			emptyList := make([]BasicItem, 0)
			emptyList = append(emptyList, BasicItem{
				Title:        "can't find any projects",
				SubTitle:     "you can execute command <jetbrains add xxxx> to add projects",
				Match:        "can't find any projects3",
				AutoComplete: "can't find any projects4",
			})
			emptyResult := EmptyResult{
				BasicItems: emptyList,
			}

			emptyJson, _ := json.Marshal(emptyResult)

			emptyJsonString := string(emptyJson)
			fmt.Println(emptyJsonString)

			os.Exit(0)
		}

		result := Result{
			Items: ItemsList,
			Variables: Variables{
				JbBin: bin,
			},
		}

		j, _ := json.Marshal(result)

		jsonX := string(j)
		fmt.Println(jsonX)
	},
}
