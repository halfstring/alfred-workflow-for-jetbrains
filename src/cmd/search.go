package cmd

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
	Arg          string    `json:"arg"`
	Title        string    `json:"title"`
	SubTitle     string    `json:"subtitle"`
	Match        string    `json:"match"`
	Valid        bool      `json:"valid"`
	AutoComplete string    `json:"autocomplete"`
	Type         string    `json:"type"`
	Variables    Variables `json:"variables"`
}

type Variables struct {
	JbBin string `json:"jb_bin"`
}

func init() {
	rootCmd.AddCommand(searchCommand)

	searchCommand.Flags().String("keyword", "", "Please input your keyword to continue.")
}

var searchCommand = &cobra.Command{
	Use:   "search",
	Short: "filter project",
	Run: func(cmd *cobra.Command, args []string) {
		keyword, _ := cmd.Flags().GetString("keyword")
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
			if len(p) != 2 {
				continue
			}

			if keyword == "" || strings.Contains(strings.ToLower(p[0]), strings.ToLower(keyword)) {
				ItemsList = append(ItemsList, Item{
					Arg:          p[1],
					Title:        p[0],
					SubTitle:     p[1],
					Match:        p[0],
					Valid:        true,
					AutoComplete: p[0],
					Type:         "default",
					Variables: Variables{
						JbBin: bin,
					},
				})

			}
		}

		if len(ItemsList) < 1 {
			type EmptyItem struct {
				Title        string `json:"title"`
				SubTitle     string `json:"subtitle"`
				Match        string `json:"match"`
				AutoComplete string `json:"autocomplete"`
			}

			type EmptyResult struct {
				EmptyItems []EmptyItem `json:"items"`
			}

			emptyList := make([]EmptyItem, 0)
			emptyList = append(emptyList, EmptyItem{
				Title:        "can't find any projects",
				SubTitle:     "you can execute command <jetbrains add xxxx> to add projects",
				Match:        "can't find any projects3",
				AutoComplete: "can't find any projects4",
			})
			emptyResult := EmptyResult{
				EmptyItems: emptyList,
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
