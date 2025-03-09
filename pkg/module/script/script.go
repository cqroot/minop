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

package script

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cqroot/minop/pkg/module"
	"github.com/cqroot/minop/pkg/module/command"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/minop/pkg/utils"
	"github.com/fatih/color"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type Module struct {
	r      *remote.Remote
	script string
}

func New(r *remote.Remote, argMap map[string]string) (*Module, error) {
	c := Module{
		r: r,
	}

	err := module.ValidateArgs(argMap, Doc.Args())
	if err != nil {
		return nil, err
	}

	c.script = argMap["script"]
	return &c, nil
}

func (m *Module) Run(resultsCh chan string) error {
	resultsCh <- color.New(color.FgYellow).Sprintf("[%s] 🟢 [%s@%s] Run local script %s on the remote node.", utils.TimeString(), m.r.Username, m.r.Hostname, m.script)
	remoteFilename := fmt.Sprintf("/tmp/minop.tmp_%s", RandomString(10))

	err := m.r.UploadFile(m.script, remoteFilename)
	if err != nil {
		resultsCh <- color.New(color.FgRed).Sprintf("[%s] 📤 [%s@%s] Failed to transfer file from %s to %s. Reason: %s.", utils.TimeString(), m.r.Username, m.r.Hostname, m.script, remoteFilename, err.Error())
		return err
	}
	defer func() {
		removeCommand, err := command.New(m.r, map[string]string{
			"command":           "rm -f " + remoteFilename,
			"print_intro":       "false",
			"print_exit_status": "false",
			"print_stdout":      "false",
			"print_stderr":      "false",
		})
		if err != nil {
			return
		}
		_ = removeCommand.Run(resultsCh)
	}()

	scriptCommand, err := command.New(m.r, map[string]string{
		"command":     fmt.Sprintf("/bin/bash %s", remoteFilename),
		"print_intro": "false",
	})
	if err != nil {
		return err
	}

	return scriptCommand.Run(resultsCh)
}
