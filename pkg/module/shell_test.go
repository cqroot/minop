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

func TestNewShell(t *testing.T) {
	tests := []struct {
		name    string
		input   module.Task
		wantErr bool
	}{
		{
			name: "valid shell input",
			input: module.Task{
				Shell: "echo hello",
			},
			wantErr: false,
		},
		{
			name: "empty shell",
			input: module.Task{
				Shell: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := module.NewShell(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, m)
			} else {
				require.NoError(t, err)
				require.NotNil(t, m)
			}
		})
	}
}

func TestShell_Name(t *testing.T) {
	m, err := module.NewShell(module.Task{Shell: "echo hello"})
	require.NoError(t, err)

	m.SetName("custom name")
	require.Equal(t, "custom name", m.Name())
}

func TestShell_DefaultName(t *testing.T) {
	m, err := module.NewShell(module.Task{Shell: "ls -la"})
	require.NoError(t, err)

	require.Equal(t, "[shell] ls -la", m.DefaultName())
}

func TestShell_Role(t *testing.T) {
	m, err := module.NewShell(module.Task{Shell: "echo hello"})
	require.NoError(t, err)

	m.SetRole("web")
	require.Equal(t, "web", m.Role())
}
