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
	"errors"
	"fmt"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/remote"
)

const (
	OpTypeShell = "shell"
	OpTypeCopy  = "copy"
	OpTypeLocal = "local"
)

type Input struct {
	Name string `yaml:"name"`
	Role string `yaml:"role"`
	Type string `yaml:"type"`

	Shell string `yaml:"shell"`
	Local string `yaml:"local"`

	Copy   string `yaml:"copy"`
	To     string `yaml:"to"`
	Backup bool   `yaml:"backup"`
}

type Operation interface {
	baseOperation
	Execute(r *remote.Remote) (*gtypes.OrderedMap[string, string], error)
	DefaultName() string
}

var (
	ErrInvalidOperation = errors.New("invalid operation")
	ErrInvalidOpType   = errors.New("invalid operation type")
)

func MakeErrInvalidOperation(in Input) error {
	return fmt.Errorf("%w: %+v", ErrInvalidOperation, in)
}

func GetOperation(in Input) (Operation, error) {
	opType := in.Type

	if opType == "" {
		if in.Shell != "" {
			opType = OpTypeShell
		} else if in.Copy != "" {
			opType = OpTypeCopy
		} else if in.Local != "" {
			opType = OpTypeLocal
		}
	}

	switch opType {
	case OpTypeShell:
		return NewOpShell(in)
	case OpTypeCopy:
		return NewOpCopy(in)
	case OpTypeLocal:
		return NewOpLocal(in)
	default:
		return nil, fmt.Errorf("%w: %q (expected %q, %q or %q)", ErrInvalidOpType, opType, OpTypeShell, OpTypeCopy, OpTypeLocal)
	}
}
