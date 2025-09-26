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
	"errors"
	"strconv"

	"github.com/cqroot/gtypes/orderedmap"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/minop/pkg/utils/maputils"
)

var ErrEmptyCommand = errors.New("empty command")

type CmdAct struct {
	Command string
}

func New(actCtx map[string]string) (*CmdAct, error) {
	act := CmdAct{}
	if err := act.Validate(actCtx); err != nil {
		return nil, err
	}
	return &act, nil
}

func (act *CmdAct) Validate(actCtx map[string]string) error {
	cmd, err := maputils.GetString(actCtx, "command")
	if err != nil {
		return err
	}
	if cmd == "" {
		return ErrEmptyCommand
	}
	act.Command = cmd
	return nil
}

func (act *CmdAct) Execute(r *remote.Remote, logger *log.Logger) (*orderedmap.OrderedMap[string, string], error) {
	exitStatus, stdout, stderr, err := r.ExecuteCommand(act.Command)
	if err != nil {
		return nil, err
	}

	res := orderedmap.New[string, string]()
	res.Put("ExitStatus", strconv.Itoa(exitStatus))
	res.Put("Stdout", stdout)
	res.Put("Stderr", stderr)
	return res, nil
}
