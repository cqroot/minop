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
	"github.com/cqroot/minop/pkg/module"
)

type ModuleDoc struct {
	desc string
	args []module.Arg
}

var Doc = ModuleDoc{
	desc: "Runs a local script on a remote node after transferring it.",
	args: []module.Arg{
		{Name: "script", Desc: "Path to the local script to run followed by optional arguments.", Type: module.ArgTypeString, Optional: false, Default: ""},
	},
}

func (md ModuleDoc) Desc() string {
	return md.desc
}

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (md ModuleDoc) Args() []module.Arg {
	return md.args
}
