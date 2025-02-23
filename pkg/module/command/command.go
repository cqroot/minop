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

package command

import (
	"errors"
	"fmt"

	"github.com/cqroot/minop/pkg/module"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/minop/pkg/utils"
	"github.com/fatih/color"
)

var ErrExitStatus = errors.New("exit status not zero")

type Module struct {
	r   *remote.Remote
	cmd string
}

func New(r *remote.Remote, argMap map[string]string) (*Module, error) {
	c := Module{
		r: r,
	}

	err := module.ValidateArgs(argMap, Doc.Args())
	if err != nil {
		return nil, err
	}

	c.cmd = argMap["command"]
	return &c, nil
}

func formattedString(fg color.Attribute, emoji string, r *remote.Remote, msg string) string {
	return color.New(fg).Sprintf("[%s] %s [%s@%s] %s", utils.TimeString(), emoji, r.Username, r.Hostname, msg)
}

func (m *Module) Run(resultsCh chan string) error {
	resultsCh <- formattedString(color.FgYellow, "🟢", m.r, fmt.Sprintf("Command: %s", m.cmd))

	exitStatus, stdout, stderr, err := m.r.ExecuteCommand(m.cmd)
	if err != nil {
		resultsCh <- fmt.Sprintf("%s %s",
			formattedString(color.FgRed, "❗", m.r, "Error:"), err.Error())
		return err
	}

	err = nil
	if exitStatus == 0 {
		resultsCh <- fmt.Sprintf("%s %d", formattedString(color.FgGreen, "✅", m.r, "Exit Status:"), exitStatus)
	} else {
		resultsCh <- fmt.Sprintf("%s %d", formattedString(color.FgRed, "❎", m.r, "Exit Status:"), exitStatus)
		err = fmt.Errorf("%w: %d", ErrExitStatus, exitStatus)
	}

	if stdout != "" {
		resultsCh <- fmt.Sprintf("%s\n%s", formattedString(color.FgCyan, "📄", m.r, "Stdout:"), stdout)
	}

	if stderr != "" {
		resultsCh <- fmt.Sprintf("%s\n%s", formattedString(color.FgRed, "🚨 ", m.r, "Stderr:"), stderr)
	}

	return err
}
