package root

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "rlt [command]",
		Short:        "Registry Load Tester",
		Long:         "A tool to test the load on a registry by running multiple instances.",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		authCmd(),
		pullCmd(),
	)
	return cmd
}
