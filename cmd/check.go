/*
Copyright (C) 2025 Keith Chu <cqroot@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"

	"github.com/cqroot/minop/pkg/executor"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func RunCheckCmd(cmd *cobra.Command, args []string) {
	e := executor.New(
		executor.WithVerboseLevel(flagVerboseLevel),
		executor.WithMaxProcs(flagMaxProcs))

	_, ops, err := e.LoadConfig(flagConfigFile)
	CheckErr(err)

	fmt.Println()
	for _, op := range ops {
		fmt.Printf("  %s %s\n", color.HiBlueString("•"), op.DefaultName())
	}
}

func NewCheckCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "check",
		Short: "Check and validate task file",
		Long:  "Load and validate the specified task file, then print the task information.",
		Run:   RunCheckCmd,
	}

	return &c
}
