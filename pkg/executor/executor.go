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

package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/logs"
	"github.com/cqroot/minop/pkg/module"
	"github.com/cqroot/minop/pkg/remote"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Executor struct {
	optVerboseLevel int
	optMaxProcs     int
}

func New(opts ...Option) *Executor {
	e := Executor{
		optVerboseLevel: 0,
		optMaxProcs:     constants.DefaultMaxProcs,
	}

	for _, opt := range opts {
		opt(&e)
	}

	return &e
}

type execResult struct {
	h   remote.Host
	res *gtypes.OrderedMap[string, string]
}

func (e Executor) ExecuteOnHosts(
	outputPrefix string,
	hostGroup map[string][]remote.Host,
	pool *remote.HostPool,
	m module.Module,
) error {
	if local, ok := m.(*module.Local); ok {
		local.SetPrefix(outputPrefix)
		res, err := local.Execute(nil)
		if err != nil {
			return err
		}
		if res != nil {
			_ = res.ForEach(func(key, val string) error {
				printKeyValue(outputPrefix, key, val, e.optVerboseLevel)
				return nil
			})
		}
		return nil
	}

	totalHosts := 0
	for _, hosts := range hostGroup {
		totalHosts += len(hosts)
	}
	logs.Logger().Debug().
		Str("module", m.Name()).
		Str("role", m.Role()).
		Int("max_procs", e.optMaxProcs).
		Int("total_hosts", totalHosts).
		Msg("ExecuteOnHosts start")

	results := make(chan execResult)

	printDone := make(chan struct{})
	go func() {
		defer close(printDone)
		for res := range results {
			printHostResult(outputPrefix, res.h, res.res, e.optVerboseLevel)
		}
	}()

	sem := semaphore.NewWeighted(int64(e.optMaxProcs))
	g, ctx := errgroup.WithContext(context.Background())

	var (
		firstErr   error
		firstErrMu sync.Mutex
	)
	recordFirstErr := func(hostStr string, err error) {
		firstErrMu.Lock()
		if firstErr == nil {
			firstErr = fmt.Errorf("host %s: %w", hostStr, err)
		}
		firstErrMu.Unlock()
	}

	for role, hosts := range hostGroup {
		if m.Role() != constants.RoleAll && m.Role() != role {
			logs.Logger().Debug().
				Str("module", m.Name()).
				Str("module_role", m.Role()).
				Str("group_role", role).
				Msg("skip host group: module role does not match")
			continue
		}

		logs.Logger().Debug().
			Str("module", m.Name()).
			Str("role", role).
			Int("host_count", len(hosts)).
			Msg("dispatching module to host group")

		for _, h := range hosts {
			hostStr := remote.HostStr(h, "")

			if err := sem.Acquire(ctx, 1); err != nil {
				logs.Logger().Error().
					Err(err).
					Str("module", m.Name()).
					Str("host", hostStr).
					Msg("semaphore acquire failed")
				if ctx.Err() != nil {
					if firstErr != nil {
						logs.Logger().Error().
							Err(firstErr).
							Str("module", m.Name()).
							Str("host", hostStr).
							Msg("context canceled due to an earlier host failure")
						return firstErr
					}
					return fmt.Errorf("context canceled while waiting for host %s: %w", hostStr, err)
				}
				continue
			}

			r, err := pool.GetRemote(h)
			if err != nil {
				sem.Release(1)
				logs.Logger().Error().
					Err(err).
					Str("module", m.Name()).
					Str("host", hostStr).
					Msg("get remote connection failed")
				return fmt.Errorf("get remote for host %s: %w", hostStr, err)
			}

			currHost := h
			g.Go(func() error {
				defer sem.Release(1)

				logs.Logger().Debug().
					Str("module", m.Name()).
					Str("host", hostStr).
					Msg("executing module on host")

				res, err := m.Execute(r)
				if err != nil {
					recordFirstErr(hostStr, err)
					logs.Logger().Error().
						Err(err).
						Str("module", m.Name()).
						Str("host", hostStr).
						Msg("module execution failed")
					return err
				}

				results <- execResult{
					h:   currHost,
					res: res,
				}
				return nil
			})
		}
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- g.Wait()
		close(results)
	}()

	<-printDone
	if err := <-errCh; err != nil {
		if firstErr != nil {
			return firstErr
		}
		return err
	}
	return nil
}

func (e Executor) ExecuteModules(
	outputPrefix string,
	hostGroup map[string][]remote.Host,
	modules []module.Module,
) error {
	termWidth := getTerminalWidth()
	pool := remote.NewHostPool()

	for _, m := range modules {
		printTaskHeader(m.Name(), termWidth)

		if err := e.ExecuteOnHosts(outputPrefix, hostGroup, pool, m); err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}
