package dupester

import (
	"fmt"
	"os"
)

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

import (
	"github.com/Eagerod/dupester/pkg/dupester"
)

var cfgFile string
var dupesterClient *dupester.Dupester

func init() {
	cobra.OnInitialize(initConfig, initDupester)
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dupester",
		Short: "command line tool for finding duplicates, or something",
		Args:  cobra.ArbitraryArgs,
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dupester.yaml)")

	return cmd
}

func Run() error {
	cmd := rootCmd()

	cmd.AddCommand(addCmd())
	cmd.AddCommand(checkCmd())

	return cmd.Execute()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".dupester" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dupester")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initDupester() {
	var tikaString string
	viper.UnmarshalKey("tika", &tikaString)

	var elasticsearchString string
	viper.UnmarshalKey("elasticsearch", &elasticsearchString)

	ds, err := dupester.NewDupester(tikaString, elasticsearchString)
	if err != nil {
		panic(err)
	}

	dupesterClient = ds
}
