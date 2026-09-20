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
	"errors"
	"testing"

	"github.com/cqroot/minop/pkg/module"
	"github.com/stretchr/testify/require"
)

func TestGetModule(t *testing.T) {
	tests := []struct {
		name     string
		input    module.Task
		wantType string
		wantErr  error
	}{
		{
			name: "shell inferred from shell field",
			input: module.Task{
				Shell: "echo hello",
			},
			wantType: "*module.Shell",
			wantErr:  nil,
		},
		{
			name: "local inferred from local field",
			input: module.Task{
				Local: "echo hello",
			},
			wantType: "*module.Local",
			wantErr:  nil,
		},
		{
			name: "copy inferred from copy body",
			input: module.Task{
				Copy: &module.CopySpec{
					Src:  "/local/path",
					Dest: "/remote/path",
				},
			},
			wantType: "*module.Copy",
			wantErr:  nil,
		},
		{
			name:     "empty input",
			input:    module.Task{},
			wantType: "",
			wantErr:  module.ErrInvalidModuleSpec,
		},
		{
			name: "shell with surrounding whitespace counts as shell",
			input: module.Task{
				Shell: "  echo hello  ",
			},
			wantType: "*module.Shell",
			wantErr:  nil,
		},
		{
			name: "multiple actions is rejected",
			input: module.Task{
				Shell: "echo hello",
				Local: "echo world",
			},
			wantType: "",
			wantErr:  module.ErrInvalidModuleSpec,
		},
		{
			name: "shell and copy is rejected",
			input: module.Task{
				Shell: "echo hello",
				Copy:  &module.CopySpec{Src: "/a", Dest: "/b"},
			},
			wantType: "",
			wantErr:  module.ErrInvalidModuleSpec,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := module.GetModule(tt.input)

			if tt.wantErr != nil {
				require.True(t, errors.Is(err, tt.wantErr))
				require.Nil(t, m)
			} else {
				require.NoError(t, err)
				require.NotNil(t, m)

				switch tt.wantType {
				case "*module.Shell":
					_, ok := m.(*module.Shell)
					require.True(t, ok, "expected Shell")
				case "*module.Copy":
					_, ok := m.(*module.Copy)
					require.True(t, ok, "expected Copy")
				case "*module.Local":
					_, ok := m.(*module.Local)
					require.True(t, ok, "expected Local")
				}
			}
		})
	}
}
