package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"jetbrains-workflow/golang/cmd/jetbrains"
	"jetbrains-workflow/golang/cmd/languages"
	"os"
)

var (
	// Used for flags.
	cfgFile string

	RootCmd = &cobra.Command{
		Use:   "root",
		Short: "Jetbrains项目管理工具",
		Long:  `Jetbrains项目管理工具`,
	}
)

// Execute executes the root command.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	//rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	//rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	RootCmd.PersistentFlags().Bool("debug", true, "Debug Mod")
	RootCmd.PersistentFlags().String("bin", "/usr/local/bin/pstorm", "Application Path")
	RootCmd.PersistentFlags().String("data-dir", os.Getenv("HOME"), "Storage Data Dir")

	RootCmd.AddCommand(jetbrains.AddCmd)
	RootCmd.AddCommand(jetbrains.ListCommand)
	RootCmd.AddCommand(jetbrains.RemoveCommand)
	RootCmd.AddCommand(jetbrains.SearchCommand)

	RootCmd.AddCommand(languages.GoogleCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
