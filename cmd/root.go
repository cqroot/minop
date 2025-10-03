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

	"github.com/cqroot/minop/pkg/cli"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/version"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	logger = log.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05 Mon"}).
		Level(zerolog.ErrorLevel)
	flagTaskFile     string
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
	logger.Debug().
		Str("task_file", flagTaskFile).
		Int("max_procs", flagMaxProcs).
		Int("verbose_level", flagVerboseLevel).
		Msg("run root command")

	if flagTaskFile == "" {
		c := cli.New(logger, cli.WithMaxProcs(flagMaxProcs), cli.WithVerboseLeve(flagVerboseLevel))
		err := c.Run()
		CheckErr(err)
		return
	}

	e := executor.New(logger,
		executor.WithVerboseLeve(flagVerboseLevel),
		executor.WithMaxProcs(flagMaxProcs))

	ops, err := e.LoadOperations(flagTaskFile, logger)
	CheckErr(err)

	err = e.ExecuteOperations(ops)
	CheckErr(err)
}

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "minop",
		Short: "MINOP is a simple tool for remote task orchestration and batch execution.",
		Long:  "MINOP is a simple tool for remote task orchestration and batch execution.",
		Run:   RunRootCmd,
	}
	rootCmd.Flags().StringVarP(&flagTaskFile, "task", "t", "", "Specify task file")
	rootCmd.PersistentFlags().IntVarP(&flagMaxProcs, "max-procs", "p", 1, "Maximum number of tasks to execute simultaneously (default 1)")
	rootCmd.Flags().CountVarP(&flagVerboseLevel, "verbose", "v", "Increase output verbosity. Use multiple v's for more detail, e.g., -v, -vv (default 0)")

	rootCmd.AddCommand(NewHostCmd())
	rootCmd.AddCommand(NewCliCmd())
	rootCmd.Version = version.Get().String()
	return &rootCmd
}

func Execute() {
	CheckErr(NewRootCmd().Execute())
}
