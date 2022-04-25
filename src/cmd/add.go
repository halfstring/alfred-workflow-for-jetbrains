package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"jetbrains-workflow/golang/tools"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().String("name", "", "Project Name.")
	addCmd.Flags().String("path", "", "Project Path.")
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add project",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		path, _ := cmd.Flags().GetString("path")
		if name == "" {
			println("name不能为空")
			os.Exit(1)
		}

		if path == "" {
			path, _ = os.Getwd()
		}

		dataDir, _ := tools.GetDataDir(cmd)
		filePath := dataDir + "/projects.list"

		_, b := tools.IsFile(filePath)

		if !b {
			f, _ := os.Create(filePath)
			defer f.Close()
		}

		fil, err := os.OpenFile(filePath, os.O_WRONLY, 0666)
		if err != nil {
			defer fil.Close()

			fmt.Println("文件打开失败", err)
			os.Exit(1)
		}

		rfp, _ := os.Open(filePath)
		defer rfp.Close()

		buf := bufio.NewScanner(rfp)
		projects := make(map[string]string)
		for {
			if !buf.Scan() {
				break //文件读完了,退出for
			}
			line := buf.Text() //获取每一行

			p := strings.Split(line, "\t")
			if len(p) != 2 {
				continue
			}
			projects[p[0]] = p[1]
		}

		if _, exist := projects[name]; exist {
			println("该project已存在")
			os.Exit(1)
		}

		n, _ := fil.Seek(0, os.SEEK_END)
		_, err = fil.WriteAt([]byte(name+"\t"+path+"\n"), n)

		fmt.Println(name + " is added successfully. ")
	},
}
