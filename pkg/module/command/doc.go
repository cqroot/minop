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

import "github.com/cqroot/minop/pkg/module"

type ModuleDoc struct {
	desc string
	args []module.Arg
}

var Doc = ModuleDoc{
	desc: "Execute commands on targets.",
	args: []module.Arg{
		{Name: "command", Desc: "The command to run.", Type: module.ArgTypeString, Optional: false, Default: ""},
		{Name: "print_intro", Desc: "Print introduction.", Type: module.ArgTypeBoolean, Optional: true, Default: "true"},
		{Name: "print_exit_status", Desc: "Print exit status.", Type: module.ArgTypeBoolean, Optional: true, Default: "true"},
		{Name: "print_stdout", Desc: "Print stdout.", Type: module.ArgTypeBoolean, Optional: true, Default: "true"},
		{Name: "print_stderr", Desc: "Print stderr.", Type: module.ArgTypeBoolean, Optional: true, Default: "true"},
	},
}

func (md ModuleDoc) Desc() string {
	return md.desc
}

func (md ModuleDoc) Args() []module.Arg {
	return md.args
}
