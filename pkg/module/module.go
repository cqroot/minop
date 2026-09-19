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
// A task declares exactly one action: a "shell" string, a "local"
// string, or a nested "copy" object. Which key is present determines
// the operation type — mirroring ansible's task module convention.
// Name should be set so that `minop task` shows a human-readable
// label; Role defaults to "all".
type Task struct {
	Name string `yaml:"name"`
	Role string `yaml:"role"`

	Shell string    `yaml:"shell,omitempty"`
	Local string    `yaml:"local,omitempty"`
	Copy  *CopySpec `yaml:"copy,omitempty"`
}

type Module interface {
	// Name and Role identify the module inside the task list and
	// determine which host group it targets. SetName / SetRole are
	// setters so the YAML loader and other code paths can populate
	// them after construction.
	Name() string
	SetName(name string)
	Role() string
	SetRole(role string)

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
// the per-op constructor uses ErrInvalidModule to refuse a
// half-filled action (e.g. copy with src but no dest).
var ErrInvalidModule = errors.New("invalid module fields")

func MakeErrInvalidModule(in Task) error {
	// Report the task name and the missing field for the action the
	// user declared, instead of dumping the whole Task struct — which
	// adds noise and would risk leaking any future sensitive fields
	// (passwords, tokens, etc.). GetModule has already selected the
	// action key, so we only describe what is missing within it.
	var missing []string
	switch {
	case in.Shell != "":
		// shell is a single string; if it were missing, GetModule
		// would have routed elsewhere, so we never get here.
	case in.Local != "":
		// local is a single string; same reasoning.
	case in.Copy != nil:
		if in.Copy.Src == "" {
			missing = append(missing, "copy.src")
		}
		if in.Copy.Dest == "" {
			missing = append(missing, "copy.dest")
		}
	default:
		// No action key was set; GetModule would have already
		// errored, but be defensive.
		missing = append(missing, "shell/local/copy")
	}

	name := in.Name
	if name == "" {
		name = "<unnamed>"
	}
	if len(missing) == 0 {
		return fmt.Errorf("%w: task %q", ErrInvalidModule, name)
	}
	return fmt.Errorf("%w: task %q missing %s", ErrInvalidModule, name, strings.Join(missing, ", "))
}

func GetModule(in Task) (Module, error) {
	hasShell := strings.TrimSpace(in.Shell) != ""
	hasLocal := strings.TrimSpace(in.Local) != ""
	hasCopy := in.Copy != nil

	count := 0
	if hasShell {
		count++
	}
	if hasLocal {
		count++
	}
	if hasCopy {
		count++
	}

	switch count {
	case 0:
		return nil, fmt.Errorf("%w: task must define exactly one of shell, local, or copy",
			ErrInvalidModuleSpec)
	case 1:
		// fall through to the action-specific dispatcher below
	default:
		return nil, fmt.Errorf("%w: task defines multiple actions; use only one of shell, local, or copy",
			ErrInvalidModuleSpec)
	}

	switch {
	case hasShell:
		return NewShell(in)
	case hasLocal:
		return NewLocal(in)
	case hasCopy:
		return NewCopy(in)
	}
	// Unreachable: count == 1 guarantees one of the three branches
	// above matched.
	return nil, fmt.Errorf("%w: unreachable", ErrInvalidModuleSpec)
}
