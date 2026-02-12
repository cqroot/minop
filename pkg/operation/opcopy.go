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
	"fmt"
	"os"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/logs"
	"github.com/cqroot/minop/pkg/remote"
)

type OpCopy struct {
	baseOperationImpl
	copy   string
	to     string
	backup bool
	mode   string
	owner  string
}

func NewOpCopy(in Input) (*OpCopy, error) {
	if in.To == "" {
		return nil, MakeErrInvalidOperation(in)
	}
	return &OpCopy{
		copy:   in.Copy,
		to:     in.To,
		backup: in.Backup,
		mode:   in.Mode,
		owner:  in.Owner,
	}, nil
}

func (op OpCopy) Execute(r *remote.Remote) (*gtypes.OrderedMap[string, string], error) {
	if op.backup {
		logs.Logger().Debug().Str("Dst", op.to).Msg("backup file")
		ret, stdout, stderr, err := r.ExecuteCommand(fmt.Sprintf(
			"if [ ! -e '%[1]s.minop_bak' ]; then [ -f '%[1]s' ] && cp -a -- '%[1]s' '%[1]s.minop_bak'; else exit 0; fi", op.to))
		if err != nil {
			logs.Logger().Err(err).Msg("failed to back up source file")
			return nil, err
		}
		if ret != 0 {
			err := fmt.Errorf("command ret: %d, out: %s, err: %s", ret, stdout, stderr)
			logs.Logger().Err(err).Msg("failed to back up source file")
			return nil, err
		}
	}

	fileInfo, err := os.Lstat(op.copy)
	if err != nil {
		logs.Logger().Err(err).Msg("")
		return nil, err
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		err = fmt.Errorf("%s is a symbolic link", op.copy)
		logs.Logger().Err(err).Msg("")
		return nil, err
	} else if fileInfo.IsDir() {
		err = r.UploadDir(op.copy, op.to)
	} else {
		err = r.UploadFile(op.copy, op.to)
	}

	if err != nil {
		return nil, err
	}

	res := gtypes.NewOrderedMap[string, string]()
	res.Put("Result", fmt.Sprintf("%s -> %s", op.copy, op.to))
	return res, nil
}
