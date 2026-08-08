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
	"errors"
	"testing"

	"github.com/cqroot/minop/pkg/operation"
	"github.com/stretchr/testify/require"
)

func TestGetOperation(t *testing.T) {
	tests := []struct {
		name    string
		input   operation.Input
		wantType string
		wantErr  error
	}{
		{
			name: "shell type from explicit type field",
			input: operation.Input{
				Type:  "shell",
				Shell: "echo hello",
			},
			wantType: "*operation.OpShell",
			wantErr:  nil,
		},
		{
			name: "shell type from shell field",
			input: operation.Input{
				Shell: "echo hello",
			},
			wantType: "*operation.OpShell",
			wantErr:  nil,
		},
		{
			name: "copy type from explicit type field",
			input: operation.Input{
				Type: "copy",
				Copy: "/local/path",
				To:   "/remote/path",
			},
			wantType: "*operation.OpCopy",
			wantErr:  nil,
		},
		{
			name: "copy type from copy field",
			input: operation.Input{
				Copy: "/local/path",
				To:   "/remote/path",
			},
			wantType: "*operation.OpCopy",
			wantErr:  nil,
		},
		{
			name: "local type from explicit type field",
			input: operation.Input{
				Type:  "local",
				Local: "echo hello",
			},
			wantType: "*operation.OpLocal",
			wantErr:  nil,
		},
		{
			name: "local type from local field",
			input: operation.Input{
				Local: "echo hello",
			},
			wantType: "*operation.OpLocal",
			wantErr:  nil,
		},
		{
			name: "invalid type",
			input: operation.Input{
				Type: "invalid",
			},
			wantType: "",
			wantErr:  operation.ErrInvalidOpType,
		},
		{
			name:    "empty input",
			input:   operation.Input{},
			wantType: "",
			wantErr:  operation.ErrInvalidOpType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := operation.GetOperation(tt.input)

			if tt.wantErr != nil {
				require.True(t, errors.Is(err, tt.wantErr))
				require.Nil(t, op)
			} else {
				require.NoError(t, err)
				require.NotNil(t, op)

				switch tt.wantType {
				case "*operation.OpShell":
					_, ok := op.(*operation.OpShell)
					require.True(t, ok, "expected OpShell")
				case "*operation.OpCopy":
					_, ok := op.(*operation.OpCopy)
					require.True(t, ok, "expected OpCopy")
				case "*operation.OpLocal":
					_, ok := op.(*operation.OpLocal)
					require.True(t, ok, "expected OpLocal")
				}
			}
		})
	}
}
