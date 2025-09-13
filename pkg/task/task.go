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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cqroot/minop/pkg/action"
	"github.com/cqroot/minop/pkg/action/command"
	"github.com/cqroot/minop/pkg/action/file"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/utils/maputils"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type Task struct {
	actions         []action.ActionWrapper
	logger          *log.Logger
	optVerboseLevel int
}

func New(name string, logger *log.Logger, opts ...Option) (*Task, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var actCtxs []map[string]string
	err = yaml.Unmarshal(content, &actCtxs)
	if err != nil {
		return nil, err
	}

	acts := make([]action.ActionWrapper, len(actCtxs))

	for i, actCtx := range actCtxs {
		name, err := maputils.GetString(actCtx, "name")
		if err != nil {
			return nil, err
		}

		role := maputils.GetStringOrDefault(actCtx, "role", "all")

		switch name {
		case "command":
			act, err := command.New(actCtx)
			if err != nil {
				return nil, err
			}
			acts[i] = *action.New(name, role, act)
		case "file":
			act, err := file.New(actCtx)
			if err != nil {
				return nil, err
			}
			acts[i] = *action.New(name, role, act)
		}
	}

	t := Task{
		actions:         acts,
		logger:          logger,
		optVerboseLevel: 0,
	}
	for _, opt := range opts {
		opt(&t)
	}
	return &t, nil
}

func (t Task) printValue(key string, val string, prefix string) {
	if t.optVerboseLevel == 1 && strings.IndexByte(val, '\n') == -1 {
		if val != "" {
			fmt.Printf("%s%s %s\n", prefix, color.CyanString("%s:", key), val)
		}
	} else {
		color.Cyan("%s%s:\n", prefix, key)
		scanner := bufio.NewScanner(strings.NewReader(val))
		for scanner.Scan() {
			fmt.Printf("%s    %s\n", prefix, scanner.Text())
		}
	}
}

func (t Task) PrintActionResult(hostGroup map[string][]host.Host, act action.ActionWrapper, prefix string) error {
	for role, hosts := range hostGroup {
		if act.Role() != "all" && act.Role() != role {
			continue
		}
		for _, h := range hosts {
			color.HiCyan("%s%s@%s:%d", prefix, h.User, h.Address, h.Port)
			ret, err := act.Execute(h, t.logger)
			if err != nil {
				return err
			}

			if t.optVerboseLevel == 0 {
				continue
			}

			if ret != nil {
				ret.ForEach(func(key, val string) error {
					t.printValue(key, val, fmt.Sprintf("%s    ", prefix))
					return nil
				})
			}
		}
	}
	return nil
}

func (t Task) Execute() error {
	hostGroup, err := host.Read(filepath.Join(".", constants.HostFileName))
	if err != nil {
		return err
	}

	for _, act := range t.actions {
		color.HiCyan("%s:\n", act.Name())
		err := t.PrintActionResult(hostGroup, act, "    ")
		if err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}
