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

package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cqroot/minop/pkg/action"
	"github.com/cqroot/minop/pkg/action/command"
	"github.com/cqroot/minop/pkg/action/file"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/minop/pkg/utils/maputils"
	"github.com/fatih/color"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

func LoadActionsFromConfig(filename string, logger *log.Logger) ([]action.ActionWrapper, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return nil, err
	}

	var actCtxs []map[string]string
	err = yaml.Unmarshal(content, &actCtxs)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return nil, err
	}

	acts := make([]action.ActionWrapper, len(actCtxs))

	for i, actCtx := range actCtxs {
		actName, err := maputils.GetString(actCtx, "action")
		if err != nil {
			logger.Error().Err(err).Msg("")
			return nil, err
		}
		logger.Debug().Str("ActionName", actName).Msg("found an action")

		role := maputils.GetStringOrDefault(actCtx, "role", "all")
		var act action.Action
		switch actName {
		case "command":
			act, err = command.New(actCtx)
			if err != nil {
				logger.Error().Err(err).Msg("")
				return nil, err
			}
		case "file":
			act, err = file.New(actCtx)
			if err != nil {
				logger.Error().Err(err).Msg("")
				return nil, err
			}
		}
		logger.Debug().Any("Action", act).Msg("")
		acts[i] = *action.New(maputils.GetStringOrDefault(actCtx, "name", actName), role, act)
	}

	return acts, nil
}

func ExecuteActions(e *ActionExecutor, acts []action.ActionWrapper) error {
	hostGroup, err := host.Read(filepath.Join(".", constants.HostFileName))
	if err != nil {
		return err
	}

	rgs := make(map[host.Host]*remote.Remote)

	for _, act := range acts {
		termWidth := 500
		if term.IsTerminal(int(os.Stdout.Fd())) {
			width, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err == nil {
				termWidth = width
			}
		}

		fmt.Printf("%s %s %s\n",
			color.HiCyanString(act.Name()),
			color.HiBlackString(strings.Repeat("•", termWidth-len(act.Name())-2-19)),
			color.HiBlackString(time.Now().Format("2006-01-02 15:04:05")),
		)

		err := e.PrintActionResult(hostGroup, &rgs, act, "    ")
		if err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}
