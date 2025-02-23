package cmd

import (
	"github.com/cqroot/minop/pkg/manager"
	"github.com/spf13/cobra"
)

func NewDocCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "doc",
		Short: "Show module documentation",
		Run: func(cmd *cobra.Command, args []string) {
			manager.ShowModuleDocs()
		},
	}
	return &cmd
}
