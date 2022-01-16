package dupester

import (
	"github.com/spf13/cobra"
)

import (
	"github.com/Eagerod/dupester/pkg/dupester"
)

var dupesterClient *dupester.Dupester

func init() {
	ds, err := dupester.NewDupester("http://dev.internal.aleemhaji.com:9998", "http://dev.internal.aleemhaji.com:9200")
	if err != nil {
		panic(err)
	}

	dupesterClient = ds
}

func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dupester",
		Short: "command line tool for finding duplicates, or something",
		Args:  cobra.ArbitraryArgs,
	}
}

func Run() error {
	cmd := rootCmd()

	cmd.AddCommand(addCmd())
	cmd.AddCommand(checkCmd())

	return cmd.Execute()
}
