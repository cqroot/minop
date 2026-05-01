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
	"github.com/cqroot/minop/pkg/cli"
	"github.com/spf13/cobra"
)

var flagCliConfigFile string

// RunCliCmd starts the interactive CLI mode.
func RunCliCmd(cmd *cobra.Command, args []string) {
	if flagCliConfigFile == "" {
		flagCliConfigFile = "./minop.yaml"
	}

	c := cli.New(cli.WithConfigFile(flagCliConfigFile), cli.WithMaxProcs(flagMaxProcs))
	CheckErr(c.Run())
}

func NewCliCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "cli",
		Short: "Start the interactive CLI mode",
		Long:  "Start an interactive CLI to execute commands on remote hosts.",
		Run:   RunCliCmd,
	}
	c.Flags().StringVarP(&flagCliConfigFile, "config", "c", "", "Specify config file (default ./minop.yaml)")

	return &c
}
