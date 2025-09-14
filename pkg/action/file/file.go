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
	"fmt"

	"github.com/cqroot/gtypes/orderedmap"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/minop/pkg/utils/maputils"
)

type File struct {
	src    string
	dst    string
	backup bool
}

func New(actCtx map[string]string) (*File, error) {
	act := File{}
	if err := act.Validate(actCtx); err != nil {
		return nil, err
	}
	return &act, nil
}

func (act *File) Validate(actCtx map[string]string) error {
	src, err := maputils.GetString(actCtx, "src")
	if err != nil {
		return err
	}
	act.src = src

	dst, err := maputils.GetString(actCtx, "dst")
	if err != nil {
		return err
	}
	act.dst = dst

	act.backup = maputils.GetBoolOrDefault(actCtx, "backup", false)
	return nil
}

func (act *File) Execute(r *remote.Remote, logger *log.Logger) (*orderedmap.OrderedMap[string, string], error) {
	if act.backup == true {
		logger.Debug().Str("Dst", act.dst).Msg("backup file")
		r.ExecuteCommand(fmt.Sprintf(
			"[ ! -e '%[1]s.minop_bak' ] && [ -f '%[1]s' ] && cp -a -- '%[1]s' '%[1]s.minop_bak'", act.dst))
	}

	err := r.UploadFile(act.src, act.dst)
	if err != nil {
		return nil, err
	}

	res := orderedmap.New[string, string]()
	res.Put("Result", fmt.Sprintf("%s -> %s", act.src, act.dst))
	return res, nil
}
