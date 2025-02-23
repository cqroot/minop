package command

import "github.com/cqroot/minop/pkg/module"

type ModuleDoc struct {
	desc string
	args []module.Arg
}

var Doc = ModuleDoc{
	desc: "Execute commands on targets",
	args: []module.Arg{
		{Name: "command", Desc: "The command to run.", Type: module.ArgTypeString, Optional: false, Default: ""},
	},
}

func (md ModuleDoc) Desc() string {
	return md.desc
}

func (md ModuleDoc) Args() []module.Arg {
	return md.args
}
