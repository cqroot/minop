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
	"strings"
)

var ErrExitStatus = errors.New("exit status not zero")

type Command struct {
	r   *remote.Remote
	cmd string
}

func New(r *remote.Remote, argMap map[string]string) (*Command, error) {
	c := Command{
		r: r,
	}

	err := module.ValidateArgs(argMap, Doc.Args())
	if err != nil {
		return nil, err
	}

	c.cmd = argMap["command"]
	return &c, nil
}

func (c *Command) writeOutput(builder *strings.Builder, s string) {
	builder.WriteString(s)
	fmt.Println(s)
}

func (c *Command) Run() (string, error) {
	builder := strings.Builder{}

	exitStatus, stdout, stderr, err := c.r.ExecuteCommand(c.cmd)
	if err != nil {
		c.writeOutput(&builder, fmt.Sprintf("%s %s\n",
			utils.FormattedString(color.FgRed, "❗", c.r, "Error:"), err.Error()))
		return builder.String(), err
	}

	err = nil
	if exitStatus == 0 {
		c.writeOutput(&builder, fmt.Sprintf("%s %d\n", utils.FormattedString(color.FgGreen, "✅", c.r, "Exit Status:"), exitStatus))
	} else {
		c.writeOutput(&builder, fmt.Sprintf("%s %d\n", utils.FormattedString(color.FgRed, "❎", c.r, "Exit Status:"), exitStatus))
		err = fmt.Errorf("%w: %d", ErrExitStatus, exitStatus)
	}

	if stdout != "" {
		c.writeOutput(&builder, utils.FormattedString(color.FgCyan, "📄", c.r, "Stdout:"))
		c.writeOutput(&builder, fmt.Sprintf("\n%s\n", stdout))
	}

	if stderr != "" {
		c.writeOutput(&builder, utils.FormattedString(color.FgRed, "🚨 ", c.r, "Stderr:"))
		c.writeOutput(&builder, fmt.Sprintf("\n%s\n", stderr))
	}

	return builder.String(), err
}
