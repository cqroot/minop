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
	"github.com/cqroot/minop/pkg/remote"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Hosts map[string][]string `yaml:"hosts"`
	Tasks []operation.Input   `yaml:"tasks"`
}

func (e Executor) LoadConfig(filename string) (map[string][]remote.Host, []operation.Operation, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		logs.Logger().Error().Err(err).Msg("failed to read file")
		return nil, nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		logs.Logger().Error().Err(err).Msg("failed to unmarshal YAML data")
		return nil, nil, fmt.Errorf("failed to unmarshal YAML data\n%w", err)
	}

	hostGroup := make(map[string][]remote.Host)
	for role, lines := range cfg.Hosts {
		for _, line := range lines {
			h, err := remote.ParseHostLine(line)
			if err != nil {
				return nil, nil, fmt.Errorf("parse host line for role %q: %w", role, err)
			}
			hostGroup[role] = append(hostGroup[role], h)
		}
	}

	ops := make([]operation.Operation, len(cfg.Tasks))
	for idx, in := range cfg.Tasks {
		op, err := operation.GetOperation(in)
		if err != nil {
			return nil, nil, err
		}

		if in.Name != "" {
			op.SetName(in.Name)
		} else {
			op.SetName(op.DefaultName())
		}

		if in.Role != "" {
			op.SetRole(in.Role)
		} else {
			op.SetRole("all")
		}

		ops[idx] = op
	}
	return hostGroup, ops, nil
}
