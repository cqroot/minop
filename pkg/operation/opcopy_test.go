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

package operation_test

import (
	"testing"

	"github.com/cqroot/minop/pkg/operation"
	"github.com/stretchr/testify/require"
)

func TestNewOpCopy(t *testing.T) {
	tests := []struct {
		name    string
		input   operation.Input
		wantErr bool
	}{
		{
			name: "valid copy input",
			input: operation.Input{
				Copy: "/local/file",
				To:   "/remote/file",
			},
			wantErr: false,
		},
		{
			name: "empty to field",
			input: operation.Input{
				Copy: "/local/file",
				To:   "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := operation.NewOpCopy(tt.input)
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

func TestOpCopy_Name(t *testing.T) {
	op, err := operation.NewOpCopy(operation.Input{
		Copy: "/local/path",
		To:   "/remote/path",
	})
	require.NoError(t, err)

	op.SetName("custom name")
	require.Equal(t, "custom name", op.Name())
}

func TestOpCopy_DefaultName(t *testing.T) {
	op, err := operation.NewOpCopy(operation.Input{
		Copy: "/local/file.txt",
		To:   "/remote/file.txt",
	})
	require.NoError(t, err)

	require.Equal(t, "[copy] /local/file.txt => /remote/file.txt", op.DefaultName())
}

func TestOpCopy_Role(t *testing.T) {
	op, err := operation.NewOpCopy(operation.Input{
		Copy: "/local/file",
		To:   "/remote/file",
	})
	require.NoError(t, err)

	op.SetRole("storage")
	require.Equal(t, "storage", op.Role())
}
