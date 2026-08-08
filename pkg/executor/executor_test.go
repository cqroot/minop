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

package executor_test

import (
	"testing"

	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/operation"
	"github.com/stretchr/testify/require"
)

func TestExecutor_New(t *testing.T) {
	e := executor.New()
	require.NotNil(t, e)
}

func TestExecutor_Options(t *testing.T) {
	e := executor.New(
		executor.WithVerboseLevel(2),
		executor.WithMaxProcs(5),
	)
	require.NotNil(t, e)
}

func TestExecutor_LoadTaskFile_NotFound(t *testing.T) {
	e := executor.New()
	_, _, err := e.LoadTaskFile("/nonexistent/file.yaml")
	require.Error(t, err)
}

func TestOpLocal_ExecuteWithExecutor(t *testing.T) {
	e := executor.New()

	ops := []operation.Operation{
		func() operation.Operation {
			op, _ := operation.NewOpLocal(operation.Input{Local: "echo hello"})
			op.SetRole(constants.RoleAll)
			return op
		}(),
	}

	err := e.ExecuteOperations("", nil, ops)
	require.NoError(t, err)
}
