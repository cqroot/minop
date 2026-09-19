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

package module_test

import (
	"testing"

	"github.com/cqroot/minop/pkg/module"
	"github.com/stretchr/testify/require"
)

func TestNewCopy(t *testing.T) {
	tests := []struct {
		name    string
		input   module.Task
		wantErr bool
	}{
		{
			name: "valid copy input",
			input: module.Task{
				Copy: &module.CopySpec{
					Src:  "/local/file",
					Dest: "/remote/file",
				},
			},
			wantErr: false,
		},
		{
			name:    "missing copy body",
			input:   module.Task{},
			wantErr: true,
		},
		{
			name: "missing dest",
			input: module.Task{
				Copy: &module.CopySpec{Src: "/local/file"},
			},
			wantErr: true,
		},
		{
			name: "missing src",
			input: module.Task{
				Copy: &module.CopySpec{Dest: "/remote/file"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := module.NewCopy(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, op)
			} else {
				require.NoError(t, err)
				require.NotNil(t, op)
			}
		})
	}
}

func TestCopy_Name(t *testing.T) {
	op, err := module.NewCopy(module.Task{
		Copy: &module.CopySpec{Src: "/local/path", Dest: "/remote/path"},
	})
	require.NoError(t, err)

	op.SetName("custom name")
	require.Equal(t, "custom name", op.Name())
}

func TestCopy_DefaultName(t *testing.T) {
	op, err := module.NewCopy(module.Task{
		Copy: &module.CopySpec{Src: "/local/file.txt", Dest: "/remote/file.txt"},
	})
	require.NoError(t, err)

	require.Equal(t, "[copy] /local/file.txt => /remote/file.txt", op.DefaultName())
}

func TestCopy_Role(t *testing.T) {
	op, err := module.NewCopy(module.Task{
		Copy: &module.CopySpec{Src: "/local/file", Dest: "/remote/file"},
	})
	require.NoError(t, err)

	op.SetRole("storage")
	require.Equal(t, "storage", op.Role())
}
