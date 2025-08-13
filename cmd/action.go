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

	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func RunActionCommandCmd(cmd *cobra.Command, args []string) {
	hostmgr, err := host.New("./host.list")
	cobra.CheckErr(err)
	logger := log.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05 Mon"}).
		Level(zerolog.DebugLevel)

	for _, h := range hostmgr.Hosts["default"] {
		r, err := remote.New(h, logger)
		cobra.CheckErr(err)

		exitStatus, stdout, stderr, err := r.ExecuteCommand(args[0])
		cobra.CheckErr(err)

		logger.Info().
			Str("User", h.User).
			Str("Addr", h.Address).
			Int("Port", h.Port).
			Int("ExitStatus", exitStatus).
			Msg(args[0])
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
	return &actionCommandCmd
}

func NewActionCmd() *cobra.Command {
	actionCmd := cobra.Command{
		Use:   "action",
		Short: "Run the specified action.",
		Long:  "Run the specified action.",
	}

	actionCmd.AddCommand(NewActionCommandCmd())
	return &actionCmd
}
