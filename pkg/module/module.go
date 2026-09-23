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
	"strings"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/remote"
)

// ErrInvalidModuleSpec signals that a task entry does not match any
// supported module kind, either because no action was declared or
// because more than one action was declared on the same task.
var ErrInvalidModuleSpec = errors.New("invalid module spec")

// Task is the YAML schema for a single task entry in minop.yaml.
// A task declares exactly one action: a "shell" string or a nested
// "copy" object. Which key is present determines the module kind.
// Name should be set so that `minop task` shows a human-readable
// label; Group defaults to "all".
type Task struct {
	Name  string `yaml:"name"`
	Group string `yaml:"group"`

	Shell string    `yaml:"shell,omitempty"`
	Copy  *CopySpec `yaml:"copy,omitempty"`
}

type Module interface {
	// Name and Group identify the module inside the task list and
	// determine which host group it targets. SetName / SetGroup are
	// setters so the YAML loader and other code paths can populate
	// them after construction.
	Name() string
	SetName(name string)
	Group() string
	SetGroup(group string)

	// Execute runs the module against a single remote host and
	// returns a structured result map for downstream rendering.
	Execute(r *remote.Remote) (*gtypes.OrderedMap[string, string], error)

	// DefaultName returns the auto-generated label used when the
	// task entry has no explicit name in minop.yaml.
	DefaultName() string
}

// ErrInvalidModule signals that a task declared an action but is
// missing the per-action required fields (e.g. copy without src).
// This is the second-line guard: GetModule picks the action, and
// the per-module constructor uses ErrInvalidModule to refuse a
// half-filled action (e.g. copy with src but no dest).
var ErrInvalidModule = errors.New("invalid module fields")

// MakeErrInvalidModule formats an ErrInvalidModule with the given
// task name and the list of field paths that were missing. Each
// module constructor (NewShell, NewCopy) declares its own missing
// field list rather than letting this function infer it from the
// Task — GetModule has already selected the action key, so the
// dispatcher here would only ever see a single arm.
func MakeErrInvalidModule(taskName string, missing ...string) error {
	if taskName == "" {
		taskName = "<unnamed>"
	}
	if len(missing) == 0 {
		return fmt.Errorf("%w: task %q", ErrInvalidModule, taskName)
	}
	return fmt.Errorf("%w: task %q missing %s", ErrInvalidModule, taskName, strings.Join(missing, ", "))
}

func GetModule(in Task) (Module, error) {
	hasShell := strings.TrimSpace(in.Shell) != ""
	hasCopy := in.Copy != nil

	switch {
	case hasShell && hasCopy:
		return nil, fmt.Errorf("%w: use only one of shell, copy",
			ErrInvalidModuleSpec)
	case hasShell:
		return NewShell(in)
	case hasCopy:
		return NewCopy(in)
	default:
		return nil, fmt.Errorf("%w: task must define shell or copy",
			ErrInvalidModuleSpec)
	}
}
