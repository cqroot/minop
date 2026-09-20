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
	"context"
	"errors"
	"testing"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/module"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/stretchr/testify/require"
)

// failingModule is a mock that always returns the configured error from
// Execute. It satisfies module.Module structurally; the unexported
// commonModule embedded in the interface is matched by name, not by
// declaration site, so this works from a test package.
type failingModule struct {
	name string
	role string
	err  error
}

func (o *failingModule) Name() string     { return o.name }
func (o *failingModule) SetName(s string) { o.name = s }
func (o *failingModule) Role() string     { return o.role }
func (o *failingModule) SetRole(r string) { o.role = r }
func (o *failingModule) DefaultName() string {
	return "failing"
}

func (o *failingModule) Execute(_ *remote.Remote) (*gtypes.OrderedMap[string, string], error) {
	return nil, o.err
}

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

func TestExecutor_LoadHostsFile_NotFound(t *testing.T) {
	_, err := executor.LoadHostsFile("/nonexistent/hosts.yaml")
	require.Error(t, err)
}

func TestExecutor_LoadTasksFile_NotFound(t *testing.T) {
	_, err := executor.LoadTasksFile("/nonexistent/minop.yaml")
	require.Error(t, err)
}

func TestLocal_ExecuteWithExecutor(t *testing.T) {
	e := executor.New()

	modules := []module.Module{
		func() module.Module {
			m, _ := module.NewLocal(module.Task{Local: "echo hello"})
			m.SetRole(constants.RoleAll)
			return m
		}(),
	}

	err := e.ExecuteModules("", nil, modules)
	require.NoError(t, err)
}

// TestExecuteOnHosts_PreservesFirstError verifies that when one host
// fails under a maxProcs=1 semaphore, the returned error reflects the
// real host failure rather than context.Canceled. With maxProcs=1, the
// failure of the first host cancels the errgroup context, which would
// make the second host's sem.Acquire return ctx.Err() and surface a
// misleading "context canceled" error. ExecuteOnHosts records the first
// real per-host error and prefers it at both the sem.Acquire and
// g.Wait() exit points.
func TestExecuteOnHosts_PreservesFirstError(t *testing.T) {
	sentinel := errors.New("real failure from host A")

	h1, err := remote.ParseHostLine("ops:pass@127.0.0.1:9001")
	require.NoError(t, err)
	h2, err := remote.ParseHostLine("ops:pass@127.0.0.1:9002")
	require.NoError(t, err)

	pool := remote.NewHostPool()
	pool.Put(h1, remote.NewForTesting(h1))
	pool.Put(h2, remote.NewForTesting(h2))

	m := &failingModule{role: constants.RoleAll, err: sentinel}

	hostGroup := map[string][]remote.Host{"all": {h1, h2}}
	e := executor.New(executor.WithMaxProcs(1))

	err = e.ExecuteOnHosts("", hostGroup, pool, m)

	require.Error(t, err)
	require.ErrorIs(t, err, sentinel,
		"expected the original host-A failure to surface; got %v", err)
	require.NotErrorIs(t, err, context.Canceled,
		"firstErr should mask context.Canceled; got %v", err)
}
