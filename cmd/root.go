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

	"github.com/adrg/xdg"
	"github.com/cqroot/minop/pkg/cli"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/logs"
	"github.com/cqroot/minop/pkg/version"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile       string
	flagTaskFile     string
	flagMaxProcs     int
	flagVerboseLevel int
)

func CheckErr(err error) {
	if err != nil {
		logs.Logger().Err(err).Msg("")
		os.Exit(1)
	}
}

func IsDirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return info.IsDir()
}

func initConfig(cmd *cobra.Command) error {
	configDir := filepath.Join(xdg.ConfigHome, "minop")
	configFile = filepath.Join(configDir, "minop.toml")

	if !IsDirExists(configDir) {
		return nil
	}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
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

func PersistentPreRunE(cmd *cobra.Command, args []string) error {
	err := initConfig(cmd)
	if err != nil {
		return err
	}

	if flagVerboseLevel >= 2 {
		logs.SetLogger(logs.Logger().Level(zerolog.DebugLevel))
	}

	logs.Logger().Debug().
		Str("task_file", flagTaskFile).
		Int("max_procs", flagMaxProcs).
		Int("verbose_level", flagVerboseLevel).
		Str("log_level", logs.Logger().GetLevel().String()).
		Msg("run root command")

	return nil
}

func RunRootCmd(cmd *cobra.Command, args []string) {
	if flagTaskFile == "" {
		c := cli.New(cli.WithMaxProcs(flagMaxProcs), cli.WithVerboseLeve(flagVerboseLevel))
		err := c.Run()
		CheckErr(err)
		return
	}

	e := executor.New(
		executor.WithVerboseLeve(flagVerboseLevel),
		executor.WithMaxProcs(flagMaxProcs))

	ops, err := e.LoadOperations(flagTaskFile)
	CheckErr(err)

	err = e.ExecuteOperations(ops)
	CheckErr(err)
}

func NewRootCmd() *cobra.Command {
	c := cobra.Command{
		Use:               "minop",
		Short:             "MINOP is a simple tool for remote task orchestration and batch execution.",
		Long:              "MINOP is a simple tool for remote task orchestration and batch execution.",
		PersistentPreRunE: PersistentPreRunE,
		Run:               RunRootCmd,
	}
	c.PersistentFlags().StringVarP(&flagTaskFile, "task", "t", "", "Specify task file")
	c.PersistentFlags().IntVarP(&flagMaxProcs, "max-procs", "p", 1, "Maximum number of tasks to execute simultaneously (default 1)")
	c.PersistentFlags().CountVarP(&flagVerboseLevel, "verbose", "v", "Increase output verbosity. Use multiple v's for more detail, e.g., -v, -vv (default 0)")

	c.AddCommand(NewHostCmd())
	c.AddCommand(NewTaskCmd())
	c.AddCommand(NewInfoCmd())
	c.Version = version.Get().String()
	return &c
}

func Execute() {
	CheckErr(NewRootCmd().Execute())
}
