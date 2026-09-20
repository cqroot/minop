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
	"github.com/cqroot/minop/pkg/executor"
	"github.com/spf13/cobra"
)

// RunRunCmd executes every module defined in the task file on every
// host in the hosts file, in declaration order. Output for each host
// line is prefixed with four spaces so it nests visually under the
// module header.
func RunRunCmd(cmd *cobra.Command, args []string) {
	e := executor.New(
		executor.WithVerboseLevel(flagVerboseLevel),
		executor.WithMaxProcs(flagMaxProcs),
	)

	hostGroup, err := e.LoadHostsFile(flagHostsFile)
	CheckErr(err)

	modules, err := e.LoadTasksFile(flagTaskFile)
	CheckErr(err)

	err = e.ExecuteModules("    ", hostGroup, modules)
	CheckErr(err)
}

// NewRunCmd creates the run subcommand.
func NewRunCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "run",
		Short: "Execute all tasks from the task file",
		Long:  "Execute all tasks defined in the task file on the configured remote hosts.",
		Run:   RunRunCmd,
	}

	return &c
}
