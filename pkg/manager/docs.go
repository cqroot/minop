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
