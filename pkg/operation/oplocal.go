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

package operation

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/remote"
)

type OpLocal struct {
	baseOperationImpl
	local  string
	prefix string
}

func NewOpLocal(in Input) (*OpLocal, error) {
	if in.Local == "" {
		return nil, MakeErrInvalidOperation(in)
	}
	return &OpLocal{
		local: in.Local,
	}, nil
}

func (op OpLocal) DefaultName() string {
	return fmt.Sprintf("[local] %s", op.local)
}

func (op *OpLocal) SetPrefix(prefix string) {
	op.prefix = prefix
}

func (op OpLocal) Execute(r *remote.Remote) (*gtypes.OrderedMap[string, string], error) {
	cmd := exec.Command("sh", "-c", op.local)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			os.Stdout.Write([]byte(op.prefix))
			os.Stdout.Write(scanner.Bytes())
			os.Stdout.Write([]byte("\n"))
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			os.Stderr.Write([]byte(op.prefix))
			os.Stderr.Write(scanner.Bytes())
			os.Stderr.Write([]byte("\n"))
		}
	}()

	err = cmd.Wait()

	res := gtypes.NewOrderedMap[string, string]()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			res.Put("ExitStatus", strconv.Itoa(exitErr.ExitCode()))
		} else {
			res.Put("Error", err.Error())
		}
		return res, nil
	}

	res.Put("ExitStatus", "0")
	return res, nil
}
