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
	"fmt"
	"strconv"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/remote"
)

// Shell executes shell commands on remote hosts.
type Shell struct {
	commonModule
	cmd string
}

// NewShell creates a new Shell operation from the given Task.
// Returns ErrInvalidModule if Shell field is empty.
func NewShell(in Task) (*Shell, error) {
	if in.Shell == "" {
		return nil, MakeErrInvalidModule(in)
	}
	return &Shell{
		cmd: in.Shell,
	}, nil
}

// DefaultName returns the default name for shell operations.
func (op Shell) DefaultName() string {
	return fmt.Sprintf("[shell] %s", op.cmd)
}

// Execute runs the shell command on the remote host and returns the results.
func (op Shell) Execute(r *remote.Remote) (*gtypes.OrderedMap[string, string], error) {
	exitStatus, stdout, stderr, err := r.ExecuteCommand(op.cmd)
	if err != nil {
		return nil, err
	}

	res := gtypes.NewOrderedMap[string, string]()
	res.Put("ExitStatus", strconv.Itoa(exitStatus))
	res.Put("Stdout", stdout)
	res.Put("Stderr", stderr)
	return res, nil
}
