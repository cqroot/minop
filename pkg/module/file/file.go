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

package file

import (
	"errors"

	"github.com/cqroot/minop/pkg/utils"
	"github.com/fatih/color"

	"github.com/cqroot/minop/pkg/module"
	"github.com/cqroot/minop/pkg/remote"
)

var ErrExitStatus = errors.New("exit status not zero")

type Module struct {
	r   *remote.Remote
	src string
	dst string
}

func New(r *remote.Remote, argMap map[string]string) (*Module, error) {
	c := Module{
		r: r,
	}

	err := module.ValidateArgs(argMap, Doc.Args())
	if err != nil {
		return nil, err
	}

	c.src = argMap["src"]
	c.dst = argMap["dst"]
	return &c, nil
}

func (m *Module) Run(resultsCh chan string) error {
	resultsCh <- color.New(color.FgYellow).Sprintf("[%s] 🟢 [%s@%s] Transfer file from %s to %s.", utils.TimeString(), m.r.Username, m.r.Hostname, m.src, m.dst)

	err := m.r.UploadFile(m.src, m.dst)
	if err == nil {
		resultsCh <- color.New(color.FgGreen).Sprintf("[%s] 📤 [%s@%s] Successfully transferred file from %s to %s.", utils.TimeString(), m.r.Username, m.r.Hostname, m.src, m.dst)
	} else {
		resultsCh <- color.New(color.FgRed).Sprintf("[%s] 📤 [%s@%s] Failed to transfer file from %s to %s. Reason: %s.", utils.TimeString(), m.r.Username, m.r.Hostname, m.src, m.dst, err.Error())
	}
	return err
}
