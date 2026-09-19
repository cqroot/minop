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

	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/logs"
	"github.com/cqroot/minop/pkg/version"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Package-wide flag values bound by cobra to the root command's
// persistent flags. They are populated by PersistentPreRunE and
// consumed by every subcommand's Run function.
var (
	// flagTaskFile is the path to the YAML file describing the tasks
	// to execute. Defaults to ./minop.yaml when empty.
	flagTaskFile string
	// flagHostsFile is the path to the YAML file mapping host roles
	// to remote connection strings. Defaults to ./hosts.yaml when
	// empty.
	flagHostsFile string
	// flagMaxProcs is the maximum number of operations executed
	// concurrently. Defaults to constants.DefaultMaxProcs.
	flagMaxProcs int
	// flagVerboseLevel is the verbosity counter (-v, -vv, -vvv)
	// controlling the zerolog level.
	flagVerboseLevel int
)

// CheckErr logs the error and exits if err is not nil.
func CheckErr(err error) {
	if err != nil {
		logs.Logger().Err(err).Msg("")
		os.Exit(1)
	}
}

// initTask applies default paths to the empty flag values and binds
// the counted/int flags through viper so they can also be sourced
// from config files or environment variables in the future.
func initTask(cmd *cobra.Command) error {
	if flagTaskFile == "" {
		flagTaskFile = "./" + constants.DefaultTaskFile
	}
	if flagHostsFile == "" {
		flagHostsFile = "./" + constants.DefaultHostsFile
	}

	if err := viper.BindPFlag("max-procs", cmd.Flags().Lookup("max-procs")); err != nil {
		return err
	}
	flagMaxProcs = viper.GetInt("max-procs")

	if err := viper.BindPFlag("verbose", cmd.Flags().Lookup("verbose")); err != nil {
		return err
	}
	flagVerboseLevel = viper.GetInt("verbose")

	return nil
}

// PersistentPreRunE is the pre-run hook that initializes logging and task file.
func PersistentPreRunE(cmd *cobra.Command, args []string) error {
	err := initTask(cmd)
	if err != nil {
		return err
	}

	if flagVerboseLevel >= 2 {
		logs.SetLogger(logs.Logger().Level(zerolog.DebugLevel))
	}

	logs.Logger().Debug().
		Str("task_file", flagTaskFile).
		Str("hosts_file", flagHostsFile).
		Int("max_procs", flagMaxProcs).
		Int("verbose_level", flagVerboseLevel).
		Str("log_level", logs.Logger().GetLevel().String()).
		Msg("run root command")

	return nil
}

// NewRootCmd creates and returns the root cobra command. Running the
// root command without a subcommand prints help; use "minop run" to
// execute tasks.
func NewRootCmd() *cobra.Command {
	c := cobra.Command{
		Use:               "minop",
		Short:             "MINOP is a simple tool for remote task orchestration and batch execution",
		Long:              "MINOP is a simple tool for remote task orchestration and batch execution.\n\nRun 'minop run' to execute tasks defined in the task file. Use 'minop task' to list them.",
		PersistentPreRunE: PersistentPreRunE,
	}
	c.PersistentFlags().StringVarP(&flagTaskFile, "task", "t", "", "Specify task file (default ./"+constants.DefaultTaskFile+")")
	c.PersistentFlags().StringVarP(&flagHostsFile, "hosts-file", "H", "", "Specify hosts file (default ./"+constants.DefaultHostsFile+")")
	c.PersistentFlags().IntVarP(&flagMaxProcs, "max-procs", "p", constants.DefaultMaxProcs, "Maximum number of tasks to execute simultaneously")
	c.PersistentFlags().CountVarP(&flagVerboseLevel, "verbose", "v", "Increase output verbosity. Use multiple v's for more detail, e.g., -v, -vv (default 0)")

	c.AddCommand(NewRunCmd())
	c.AddCommand(NewHostCmd())
	c.AddCommand(NewTaskCmd())
	c.AddCommand(NewCliCmd())
	c.Version = version.Get().String()
	return &c
}

// Execute builds the root command tree and runs it. It is the single
// entry point invoked by main; any error returned by cobra is surfaced
// via CheckErr which logs and exits non-zero.
func Execute() {
	CheckErr(NewRootCmd().Execute())
}
