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

package module

import (
	"fmt"
	"os"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/logs"
	"github.com/cqroot/minop/pkg/remote"
)

// CopySpec is the nested body of a "copy" task. Src is the local
// path; Dest is the absolute remote path. When Backup is true,
// any pre-existing file at Dest is renamed to "Dest.minop_bak"
// before the upload, mirroring ansible's backup: yes option.
type CopySpec struct {
	Src    string `yaml:"src"`
	Dest   string `yaml:"dest"`
	Backup bool   `yaml:"backup"`
}

// Copy copies files or directories to remote hosts via SFTP.
type Copy struct {
	commonModule
	src    string
	dest   string
	backup bool
}

// NewCopy creates a new Copy operation from the given Task.
// Returns ErrInvalidModule if the copy body, src or dest is missing.
func NewCopy(in Task) (*Copy, error) {
	if in.Copy == nil || in.Copy.Src == "" || in.Copy.Dest == "" {
		return nil, MakeErrInvalidModule(in)
	}
	return &Copy{
		src:    in.Copy.Src,
		dest:   in.Copy.Dest,
		backup: in.Copy.Backup,
	}, nil
}

// DefaultName returns the default name for copy operations.
func (op Copy) DefaultName() string {
	return fmt.Sprintf("[copy] %s => %s", op.src, op.dest)
}

// Execute uploads the local file or directory to the remote host.
func (op Copy) Execute(r *remote.Remote) (*gtypes.OrderedMap[string, string], error) {
	if op.backup {
		logs.Logger().Debug().Str("Dst", op.dest).Msg("backup file")
		ret, stdout, stderr, err := r.ExecuteCommand(fmt.Sprintf(
			"if [ ! -e '%[1]s.minop_bak' ] && [ -f '%[1]s' ]; then cp -a -- '%[1]s' '%[1]s.minop_bak'; else exit 0; fi", op.dest,
		))
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

	fileInfo, err := os.Lstat(op.src)
	if err != nil {
		logs.Logger().Err(err).Str("src", op.src).Msg("stat source failed")
		return nil, err
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		err = fmt.Errorf("source %s is a symbolic link", op.src)
		logs.Logger().Err(err).Str("src", op.src).Msg("refusing to upload symlink")
		return nil, err
	}

	switch {
	case fileInfo.IsDir():
		err = r.UploadDir(op.src, op.dest)
	default:
		err = r.UploadFile(op.src, op.dest)
	}
	if err != nil {
		return nil, err
	}

	res := gtypes.NewOrderedMap[string, string]()
	res.Put("Result", fmt.Sprintf("%s -> %s", op.src, op.dest))
	return res, nil
}
