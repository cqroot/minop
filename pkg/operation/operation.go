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

type Input struct {
	Name string `yaml:"name"`
	Role string `yaml:"role"`

	Shell string `yaml:"shell"`

	Copy   string `yaml:"copy"`
	To     string `yaml:"to"`
	Backup bool   `yaml:"backup"`
	Mode   string `yaml:"mode"`
	Owner  string `yaml:"owner"`
}

type Operation interface {
	baseOperation
	Execute(r *remote.Remote) (*gtypes.OrderedMap[string, string], error)
}

var ErrInvalidOperation = errors.New("invalid operation")

func MakeErrInvalidOperation(in Input) error {
	return fmt.Errorf("%w\n%+v", ErrInvalidOperation, in)
}

func GetOperation(in Input) (Operation, error) {
	if in.Shell != "" {
		return NewOpShell(in)
	}

	if in.Copy != "" {
		return NewOpCopy(in)
	}

	return nil, MakeErrInvalidOperation(in)
}
