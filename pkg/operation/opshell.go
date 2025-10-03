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

package operation

import (
	"strconv"

	"github.com/cqroot/gtypes/orderedmap"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/rs/zerolog"
)

type OpShell struct {
	baseOperationImpl
	logger zerolog.Logger
	shell  string
}

func NewOpShell(in Input, logger zerolog.Logger) (*OpShell, error) {
	if in.Shell == "" {
		return nil, MakeErrInvalidOperation(in)
	}
	return &OpShell{
		logger: logger,
		shell:  in.Shell,
	}, nil
}

func (op OpShell) Execute(r *remote.Remote) (*orderedmap.OrderedMap[string, string], error) {
	exitStatus, stdout, stderr, err := r.ExecuteCommand(op.shell)
	if err != nil {
		return nil, err
	}

	res := orderedmap.New[string, string]()
	res.Put("ExitStatus", strconv.Itoa(exitStatus))
	res.Put("Stdout", stdout)
	res.Put("Stderr", stderr)
	return res, nil
}
