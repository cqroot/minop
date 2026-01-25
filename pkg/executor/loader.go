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

	"github.com/cqroot/minop/pkg/logs"
	"github.com/cqroot/minop/pkg/operation"
	"gopkg.in/yaml.v3"
)

func (e Executor) LoadOperations(filename string) ([]operation.Operation, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		logs.Logger().Error().Err(err).Msg("failed to read file")
		return nil, err
	}

	var ins []operation.Input
	err = yaml.Unmarshal(content, &ins)
	if err != nil {
		logs.Logger().Error().Err(err).Msg("failed to unmarshal YAML data")
		return nil, fmt.Errorf("failed to unmarshal YAML data\n%w", err)
	}

	ops := make([]operation.Operation, len(ins))
	for idx, in := range ins {
		op, err := operation.GetOperation(in)
		if err != nil {
			return nil, err
		}

		if in.Name != "" {
			op.SetName(in.Name)
		} else {
			op.SetName("Anonymous Operation")
		}

		if in.Role != "" {
			op.SetRole(in.Role)
		} else {
			op.SetRole("all")
		}

		ops[idx] = op
	}
	return ops, nil
}
