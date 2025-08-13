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
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "minop",
		Short: "Minop is a simple remote execution and deployment tool",
		Long:  "Minop is a simple remote execution and deployment tool",
	}

	rootCmd.AddCommand(NewActionCmd())
	rootCmd.AddCommand(NewHostCmd())
	return &rootCmd
}

func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}
