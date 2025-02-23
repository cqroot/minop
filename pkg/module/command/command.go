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
