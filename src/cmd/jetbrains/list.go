package jetbrains

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"jetbrains-workflow/golang/tools"
	"os"
	"strings"
)

func init() {
	//cmd.RootCmd.AddCommand(ListCommand)
}

var ListCommand = &cobra.Command{
	Use:   "project/list",
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
			if len(p) != 2 && len(p) != 3 {
				continue
			}

			if len(p) == 3 {
				fmt.Println(p[0], "\t", p[1], "\t", p[2])
			} else {
				fmt.Println(p[0], "\t", p[1])
			}

		}

		fmt.Println("Over")
	},
}
