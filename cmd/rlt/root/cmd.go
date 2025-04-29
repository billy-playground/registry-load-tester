package root

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "rlt [command]",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		pullCmd(),
	)
	return cmd
}
