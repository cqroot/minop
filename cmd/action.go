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
	"os"
	"strconv"

	"github.com/cqroot/minop/pkg/action/command"
	"github.com/cqroot/minop/pkg/action/file"
	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func NewLogger() *log.Logger {
	return log.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05 Mon"}).
		Level(zerolog.DebugLevel)
}

var (
	actCtxCommandCmd string
	actCtxFileSrc    string
	actCtxFileDst    string
)

func RunActionCommandCmd(cmd *cobra.Command, args []string) {
	hostmgr, err := host.New("./host.list")
	cobra.CheckErr(err)
	logger := NewLogger()

	for _, h := range hostmgr.Hosts["default"] {
		cmd, err := command.New(map[string]string{
			"command": actCtxCommandCmd,
		})
		cobra.CheckErr(err)

		res, err := cmd.Execute(h, logger)
		cobra.CheckErr(err)

		exitStatus, err := strconv.Atoi(res["ExitStatus"])
		cobra.CheckErr(err)
		stdout := res["Stdout"]
		stderr := res["Stderr"]

		logger.Info().
			Str("User", h.User).
			Str("Addr", h.Address).
			Int("Port", h.Port).
			Int("ExitStatus", exitStatus).
			Msg(actCtxCommandCmd)
		if stdout != "" {
			logger.Info().
				Str("User", h.User).
				Str("Addr", h.Address).
				Int("Port", h.Port).
				Msg("STDOUT")
			fmt.Println(stdout)
		}
		if stderr != "" {
			logger.Info().
				Str("User", h.User).
				Str("Addr", h.Address).
				Int("Port", h.Port).
				Msg("STDERR")
			fmt.Println(stderr)
		}
	}
}

func NewActionCommandCmd() *cobra.Command {
	actionCommandCmd := cobra.Command{
		Use:   "command",
		Short: "Execute the remote command.",
		Long:  "Execute the remote command.",
		Run:   RunActionCommandCmd,
	}
	actionCommandCmd.PersistentFlags().StringVar(&actCtxCommandCmd, "command", "", "command")
	return &actionCommandCmd
}

func RunActionFileCmd(cmd *cobra.Command, args []string) {
	hostmgr, err := host.New("./host.list")
	cobra.CheckErr(err)
	logger := NewLogger()

	for _, h := range hostmgr.Hosts["default"] {
		logger.Info().
			Str("User", h.User).
			Str("Addr", h.Address).
			Int("Port", h.Port).
			Str("Src", actCtxFileSrc).
			Str("Dst", actCtxFileDst).
			Msg("Transfer file")

		act, err := file.New(map[string]string{
			"src": actCtxFileSrc,
			"dst": actCtxFileDst,
		})
		cobra.CheckErr(err)

		_, err = act.Execute(h, logger)
		if err != nil {
			logger.Error().
				Str("User", h.User).
				Str("Addr", h.Address).
				Int("Port", h.Port).
				Str("Src", actCtxFileSrc).
				Str("Dst", actCtxFileDst).
				Err(err)
		} else {
			logger.Info().
				Str("User", h.User).
				Str("Addr", h.Address).
				Int("Port", h.Port).
				Str("Src", actCtxFileSrc).
				Str("Dst", actCtxFileDst).
				Msg("Successfully transferred file")
		}
	}
}

func NewActionFileCmd() *cobra.Command {
	actionFileCmd := cobra.Command{
		Use:   "file",
		Short: "Copy file to remote locations.",
		Long:  "Copy file to remote locations.",
		Run:   RunActionFileCmd,
	}
	actionFileCmd.PersistentFlags().StringVar(&actCtxFileSrc, "src", "", "src")
	actionFileCmd.PersistentFlags().StringVar(&actCtxFileDst, "dst", "", "dst")
	return &actionFileCmd
}

func RunActionDirCmd(cmd *cobra.Command, args []string) {
	hostmgr, err := host.New("./host.list")
	cobra.CheckErr(err)
	logger := NewLogger()

	for _, h := range hostmgr.Hosts["default"] {
		r, err := remote.New(h, logger)
		cobra.CheckErr(err)

		logger.Info().
			Str("User", h.User).
			Str("Addr", h.Address).
			Int("Port", h.Port).
			Str("Src", args[0]).
			Str("Dst", args[1]).
			Msg("Transfer dir")
		err = r.UploadDirectory(args[0], args[1])
		if err != nil {
			logger.Error().
				Str("User", h.User).
				Str("Addr", h.Address).
				Int("Port", h.Port).
				Str("Src", args[0]).
				Str("Dst", args[1]).
				Err(err)
		} else {
			logger.Info().
				Str("User", h.User).
				Str("Addr", h.Address).
				Int("Port", h.Port).
				Str("Src", args[0]).
				Str("Dst", args[1]).
				Msg("Successfully transferred dir")
		}
	}
}

func NewActionDirCmd() *cobra.Command {
	actionDirCmd := cobra.Command{
		Use:   "dir",
		Short: "Copy dir to remote locations.",
		Long:  "Copy dir to remote locations.",
		Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
		Run:   RunActionDirCmd,
	}
	return &actionDirCmd
}

func NewActionCmd() *cobra.Command {
	actionCmd := cobra.Command{
		Use:   "action",
		Short: "Run the specified action.",
		Long:  "Run the specified action.",
	}

	actionCmd.AddCommand(NewActionCommandCmd())
	actionCmd.AddCommand(NewActionFileCmd())
	actionCmd.AddCommand(NewActionDirCmd())
	return &actionCmd
}
