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

	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/logs"
	"github.com/cqroot/minop/pkg/operation"
	"github.com/cqroot/minop/pkg/remote"
	"gopkg.in/yaml.v3"
)

// hostsFileSchema is the top-level shape of hosts.yaml: a flat map from
// role name to a list of host connection strings in the
// "user:password@address:port" format.
type hostsFileSchema struct {
	Roles map[string][]string
}

// tasksFileSchema is the top-level shape of minop.yaml. Only the
// "tasks" key is recognised; hosts live in a separate file.
type tasksFileSchema struct {
	Tasks []operation.Input `yaml:"tasks"`
}

// LoadHostsFile reads a hosts file and returns a map from role name to
// parsed Host structs. The file is expected to be a flat YAML map
// keyed by role, for example:
//
//	all:
//	  - ops:PASSWORD@127.0.0.1:9001
//	main:
//	  - ops:PASSWORD@127.0.0.1:9002
//
// The outer "hosts:" wrapper that the legacy combined format used is
// intentionally not accepted; hosts.yaml is dedicated to host entries.
func (e Executor) LoadHostsFile(filename string) (map[string][]remote.Host, error) {
	logs.Logger().Debug().Str("filename", filename).Msg("loading hosts file")

	content, err := os.ReadFile(filename)
	if err != nil {
		logs.Logger().Error().Err(err).Msg("failed to read hosts file")
		return nil, err
	}

	raw := make(map[string][]string)
	if err := yaml.Unmarshal(content, &raw); err != nil {
		logs.Logger().Error().Err(err).Msg("failed to unmarshal hosts YAML")
		return nil, fmt.Errorf("failed to unmarshal hosts YAML: %w", err)
	}

	hostGroup := make(map[string][]remote.Host, len(raw))
	for role, lines := range raw {
		for _, line := range lines {
			h, err := remote.ParseHostLine(line)
			if err != nil {
				return nil, fmt.Errorf("parse host line for role %q: %w", role, err)
			}
			hostGroup[role] = append(hostGroup[role], h)
		}
	}

	return hostGroup, nil
}

// LoadTasksFile reads a task file and returns the list of operations
// to execute. Each entry's name defaults to DefaultName() and its role
// defaults to RoleAll when the corresponding YAML field is empty.
func (e Executor) LoadTasksFile(filename string) ([]operation.Operation, error) {
	logs.Logger().Debug().Str("filename", filename).Msg("loading tasks file")

	content, err := os.ReadFile(filename)
	if err != nil {
		logs.Logger().Error().Err(err).Msg("failed to read tasks file")
		return nil, err
	}

	var cfg tasksFileSchema
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		logs.Logger().Error().Err(err).Msg("failed to unmarshal tasks YAML")
		return nil, fmt.Errorf("failed to unmarshal tasks YAML: %w", err)
	}

	ops := make([]operation.Operation, len(cfg.Tasks))
	for idx, in := range cfg.Tasks {
		op, err := operation.GetOperation(in)
		if err != nil {
			return nil, err
		}

		if in.Name != "" {
			op.SetName(in.Name)
		} else {
			op.SetName(op.DefaultName())
		}

		if in.Role != "" {
			op.SetRole(in.Role)
		} else {
			op.SetRole(constants.RoleAll)
		}

		ops[idx] = op
	}
	return ops, nil
}
