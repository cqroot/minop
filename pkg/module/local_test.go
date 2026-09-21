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

func TestNewLocal(t *testing.T) {
	tests := []struct {
		name    string
		input   module.Task
		wantErr bool
	}{
		{
			name: "valid local input",
			input: module.Task{
				Local: "echo hello",
			},
			wantErr: false,
		},
		{
			name: "empty local",
			input: module.Task{
				Local: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := module.NewLocal(tt.input)
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

func TestLocal_Name(t *testing.T) {
	m, err := module.NewLocal(module.Task{Local: "echo hello"})
	require.NoError(t, err)

	m.SetName("custom name")
	require.Equal(t, "custom name", m.Name())
}

func TestLocal_DefaultName(t *testing.T) {
	m, err := module.NewLocal(module.Task{Local: "ls -la"})
	require.NoError(t, err)

	require.Equal(t, "[local] ls -la", m.DefaultName())
}

func TestLocal_Group(t *testing.T) {
	m, err := module.NewLocal(module.Task{Local: "echo hello"})
	require.NoError(t, err)

	m.SetGroup("local")
	require.Equal(t, "local", m.Group())
}

func TestLocal_SetPrefix(t *testing.T) {
	m, err := module.NewLocal(module.Task{Local: "echo hello"})
	require.NoError(t, err)

	m.SetPrefix("    ")
	require.NotNil(t, m)
}

func TestLocal_Execute(t *testing.T) {
	t.Run("successful command", func(t *testing.T) {
		m, err := module.NewLocal(module.Task{Local: "echo hello"})
		require.NoError(t, err)

		res, err := m.Execute(nil)
		require.NoError(t, err)
		require.NotNil(t, res)

		exitStatus, ok := res.Get("ExitStatus")
		require.True(t, ok)
		require.Equal(t, "0", exitStatus)
	})

	t.Run("failing command", func(t *testing.T) {
		m, err := module.NewLocal(module.Task{Local: "exit 1"})
		require.NoError(t, err)

		res, err := m.Execute(nil)
		require.NoError(t, err)
		require.NotNil(t, res)

		exitStatus, ok := res.Get("ExitStatus")
		require.True(t, ok)
		require.Equal(t, "1", exitStatus)
	})

	t.Run("command not found", func(t *testing.T) {
		m, err := module.NewLocal(module.Task{Local: "nonexistent_command_12345"})
		require.NoError(t, err)

		res, err := m.Execute(nil)
		require.NoError(t, err)
		require.NotNil(t, res)

		exitStatus, ok := res.Get("ExitStatus")
		require.True(t, ok)
		require.NotEqual(t, "0", exitStatus)
	})
}
