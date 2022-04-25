package tools

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func IsExists(path string) (os.FileInfo, bool) {
	f, err := os.Stat(path)
	return f, err == nil || os.IsExist(err)
}

func IsFile(path string) (os.FileInfo, bool) {
	f, flag := IsExists(path)
	return f, flag && !f.IsDir()
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetDataDir(cmd *cobra.Command) (dataDir string, err error) {
	dataDir, _ = cmd.Root().Flags().GetString("data-dir")

	dataDir = strings.TrimSpace(dataDir)
	dataDir = strings.TrimRight(dataDir, "/")
	if dirExist, _ := PathExists(dataDir); !dirExist {
		panic("please confirm if the path exists.")
	}

	return dataDir, nil
}
