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

package command

import (
	"fmt"
	"strconv"

	"github.com/cqroot/minop/pkg/action"
	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/remote"
)

type Command struct {
	Command string
}

func New(actCtx map[string]string) (*Command, error) {
	cmd := Command{}
	if err := cmd.Validate(actCtx); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (act *Command) Validate(actCtx map[string]string) error {
	cmd, ok := actCtx["command"]
	if !ok {
		return fmt.Errorf("%w: command", action.ErrParameterMissing)
	}
	act.Command = cmd
	return nil
}

func (act *Command) Execute(h host.Host, logger *log.Logger) (map[string]string, error) {
	r, err := remote.New(h, logger)
	if err != nil {
		return nil, err
	}

	exitStatus, stdout, stderr, err := r.ExecuteCommand(act.Command)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"ExitStatus": strconv.Itoa(exitStatus),
		"Stdout":     stdout,
		"Stderr":     stderr,
	}, nil
}
