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

package module

import (
	"errors"
	"fmt"

	"github.com/cqroot/minop/pkg/utils"
)

type ArgType string

const (
	ArgTypeString  ArgType = "String"
	ArgTypeInteger ArgType = "Integer"
	ArgTypeFloat   ArgType = "Float"
	ArgTypeBoolean ArgType = "Boolean"
)

type Arg struct {
	Name     string
	Desc     string
	Type     ArgType
	Optional bool
	Default  string
}

var (
	ErrArgNotFound     = errors.New("arg not found")
	ErrArgTypeMismatch = errors.New("arg type mismatch")
)

func ValidateArgs(argMap map[string]string, argTypes []Arg) error {
	for _, arg := range argTypes {
		// Only mandatory args require validation
		if arg.Optional {
			if argMap[arg.Name] == "" {
				argMap[arg.Name] = arg.Default
			}
			continue
		}

		val, ok := argMap[arg.Name]
		if !ok {
			return fmt.Errorf("%w: %s", ErrArgNotFound, arg.Name)
		}

		switch arg.Type {
		case ArgTypeString:
			continue

		case ArgTypeInteger:
			if !utils.StrIsInteger(val) {
				return fmt.Errorf("%w: %s %s -> %s", ErrArgTypeMismatch, arg.Name, arg.Type, val)
			}

		case ArgTypeFloat:
			if !utils.StrIsFloat(val) {
				return fmt.Errorf("%w: %s %s -> %s", ErrArgTypeMismatch, arg.Name, arg.Type, val)
			}

		case ArgTypeBoolean:
			if !utils.StrIsBoolean(val) {
				return fmt.Errorf("%w: %s %s -> %s", ErrArgTypeMismatch, arg.Name, arg.Type, val)
			}
		}
	}

	return nil
}
