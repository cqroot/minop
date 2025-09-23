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
	"os"
	"path/filepath"

	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/version"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	logger = log.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05 Mon"}).
		Level(zerolog.ErrorLevel)
	flagConfigFile   string
	flagMaxProcs     int
	flagVerboseLevel int
)

func CheckErr(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func RunRootCmd(cmd *cobra.Command, args []string) {
	if flagVerboseLevel >= 2 {
		logger = logger.Level(zerolog.DebugLevel)
	}

	e := executor.New(logger,
		executor.WithVerboseLeve(flagVerboseLevel),
		executor.WithMaxProcs(flagMaxProcs))

	acts, err := executor.LoadActionsFromConfig(flagConfigFile, logger)
	CheckErr(err)

	err = executor.ExecuteActions(e, acts)
	CheckErr(err)
}

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "minop",
		Short: "Minop is a simple remote execution and deployment tool",
		Long:  "Minop is a simple remote execution and deployment tool",
		Run:   RunRootCmd,
	}
	rootCmd.Flags().StringVarP(&flagConfigFile, "config", "c", filepath.Join(".", constants.TaskFileName), "Specify config file")
	rootCmd.Flags().IntVarP(&flagMaxProcs, "max-procs", "p", 1, "Maximum number of tasks to execute simultaneously (default 1)")
	rootCmd.Flags().CountVarP(&flagVerboseLevel, "verbose", "v", "Increase output verbosity. Use multiple v's for more detail, e.g., -v, -vv (default 0)")

	rootCmd.AddCommand(NewHostCmd())
	rootCmd.Version = version.Get().String()
	return &rootCmd
}

func Execute() {
	CheckErr(NewRootCmd().Execute())
}
