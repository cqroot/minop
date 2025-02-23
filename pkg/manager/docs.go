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

package manager

import (
	"fmt"
	"github.com/cqroot/minop/pkg/module"
	"github.com/cqroot/minop/pkg/module/command"
	"github.com/fatih/color"
)

var modules = []string{
	"command",
}

var ModuleDocMap = map[string]module.Doc{
	"command": command.Doc,
}

func ShowModuleDocs() {
	for _, name := range modules {
		fmt.Printf("%s  %s\n", color.GreenString(name), ModuleDocMap[name].Desc())
		for _, arg := range ModuleDocMap[name].Args() {
			fmt.Printf("        %s\t\t%s \t%s",
				color.CyanString(arg.Name), color.MagentaString(string(arg.Type)), arg.Desc)
			if arg.Optional {
				fmt.Printf("  %s\n", color.YellowString("(Optional, default: %s)", arg.Default))
			} else {
				fmt.Print("\n")
			}
		}
	}
}
