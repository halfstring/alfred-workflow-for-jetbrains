package jetbrains

import (
	"bufio"
	"fmt"
	"jetbrains-workflow/golang/tools"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	//cmd.RootCmd.AddCommand(RemoveCommand)

	RemoveCommand.Flags().String("name", "", "filter project name")
	//searchCommand.Flags().Bool("debug", false, "Debug Mod")
}

var RemoveCommand = &cobra.Command{
	Use:   "project/remove",
	Short: "remove a project",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		name = strings.TrimSpace(name)

		dataDir, _ := tools.GetDataDir(cmd)
		filePath := dataDir + "/projects.list"

		if name == "" {
			fmt.Println("project Name不能为空")
			os.Exit(0)
		}

		_, b := tools.IsFile(filePath)
		if !b {
			f, _ := os.Create(filePath)
			defer f.Close()

			fmt.Println("数据文件刚被创建，不需要移除项目操作。")
			os.Exit(0)
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

		//var projects = make(map[string]string)

		pStr := ""
		for {
			if !buf.Scan() {
				break
			}
			line := buf.Text()

			p := strings.Split(line, "\t")
			if len(p) != 2 && len(p) != 3 {
				continue
			}

			if strings.ToLower(p[0]) != strings.ToLower(name) {
				//projects[p[0]] = p[1]
				pStr += p[0] + "\t" + p[1]
				if len(p) == 3 {
					pStr += "\t" + p[2]
				}
				pStr += "\n"
			} else {
				fmt.Println(name, "被移除.")
			}
		}

		cmd2 := " echo \"\" > " + filePath

		o, err := exec.Command("bash", "-c", cmd2).Output()
		if err != nil {
			panic("some error found")
		}

		f2, err := os.OpenFile(filePath, os.O_WRONLY, 0666)
		defer f2.Close()

		f2.WriteString(pStr)

		fmt.Println("projects数据更新完毕。", o)
	},
}
