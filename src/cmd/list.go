package cmd

import (
	"bufio"
	"fmt"
	"jetbrains-workflow/golang/tools"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "fetch all your projects",
	Run: func(cmd *cobra.Command, args []string) {
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

			fmt.Println("文件打开失败", err)
			os.Exit(1)
		}

		fileRead, _ := os.Open(filePath)
		defer fileRead.Close()

		buf := bufio.NewScanner(fileRead)

		for {
			if !buf.Scan() {
				break
			}
			line := buf.Text()

			p := strings.Split(line, "\t")
			if len(p) != 2 {
				continue
			}

			fmt.Println(p[0], "\t", p[1])
		}

		fmt.Println("Over")
	},
}
