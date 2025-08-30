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
	"github.com/cqroot/minop/pkg/action"
	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/remote"
)

type File struct {
	LocalPath  string
	RemotePath string
}

func New(actCtx map[string]string) (*File, error) {
	act := File{}
	if err := act.Validate(actCtx); err != nil {
		return nil, err
	}
	return &act, nil
}

func (act *File) Validate(actCtx map[string]string) error {
	localPath, err := action.GetActionParam(actCtx, "local_path")
	if err != nil {
		return err
	}
	act.LocalPath = localPath

	remotePath, err := action.GetActionParam(actCtx, "remote_path")
	if err != nil {
		return err
	}
	act.RemotePath = remotePath
	return nil
}

func (act *File) Execute(h host.Host, logger *log.Logger) (map[string]string, error) {
	r, err := remote.New(h, logger)
	if err != nil {
		return nil, err
	}

	err = r.UploadFile(act.LocalPath, act.RemotePath)
	return nil, err
}
