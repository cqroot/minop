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

package task

import (
	"fmt"
	"os"

	"github.com/cqroot/minop/pkg/action"
	"github.com/cqroot/minop/pkg/action/command"
	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/utils/maputils"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type Task struct {
	Actions []action.Action
}

func New(name string) (*Task, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var actCtxs []map[string]string
	err = yaml.Unmarshal(content, &actCtxs)
	if err != nil {
		return nil, err
	}

	acts := make([]action.Action, len(actCtxs))

	for i, actCtx := range actCtxs {
		name, err := maputils.GetString(actCtx, "name")
		if err != nil {
			return nil, err
		}

		switch name {
		case "command":
			act, err := command.New(actCtx)
			if err != nil {
				return nil, err
			}
			acts[i] = act
		}
	}

	return &Task{
		Actions: acts,
	}, nil
}

func (t Task) Execute() error {
	hostGroup, err := host.Read("./host.list")
	if err != nil {
		return err
	}
	logger := log.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05 Mon"}).
		Level(zerolog.DebugLevel)

	for _, act := range t.Actions {
		for _, h := range hostGroup["default"] {
			ret, err := act.Execute(h, logger)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", ret)
		}
	}
	return nil
}
