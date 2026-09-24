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
	optPrinter      *Printer
}

func New(opts ...Option) *Executor {
	e := Executor{
		optMaxProcs: constants.DefaultMaxProcs,
	}

	for _, opt := range opts {
		opt(&e)
	}

	// Build the default Printer after the option loop so it sees the
	// final verbosity level. WithPrinter overrides this.
	if e.optPrinter == nil {
		e.optPrinter = defaultPrinter(e.optVerboseLevel)
	}

	return &e
}

type execResult struct {
	h   remote.Host
	res *gtypes.OrderedMap[string, string]
}

// firstErrorTracker captures the first non-nil per-host error so that
// downstream context-cancel-induced semaphore acquire failures can be
// hidden in favor of the original cause. See
// TestExecuteOnHosts_PreservesFirstError for the failure mode it
// guards against.
type firstErrorTracker struct {
	mu  sync.Mutex
	err error
}

func (t *firstErrorTracker) record(hostStr string, err error) {
	if err == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.err == nil {
		t.err = fmt.Errorf("host %s: %w", hostStr, err)
	}
}

func (t *firstErrorTracker) get() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.err
}

func (e Executor) ExecuteOnHosts(
	outputPrefix string,
	hostGroup map[string][]remote.Host,
	pool *remote.HostPool,
	m module.Module,
) error {
	sem := semaphore.NewWeighted(int64(e.optMaxProcs))
	tracker := &firstErrorTracker{}

	logs.Logger().Debug().
		Str("module", m.Name()).
		Str("group", m.Group()).
		Int("max_procs", e.optMaxProcs).
		Int("total_hosts", totalHostCount(hostGroup)).
		Msg("ExecuteOnHosts start")

	// The printer drains results as workers produce them. After
	// g.Wait() we close results to let the printer exit, then wait
	// on printDone so no goroutine outlives this function.
	results := make(chan execResult)
	printDone := make(chan struct{})
	go func() {
		defer close(printDone)
		for res := range results {
			e.optPrinter.PrintHostResult(outputPrefix, res.h, res.res)
		}
	}()

	g, ctx := errgroup.WithContext(context.Background())
	dispatchErr := dispatchWorkers(ctx, g, sem, pool, m, hostGroup, results, tracker)

	waitErr := g.Wait()
	close(results)
	<-printDone

	// Prefer a real per-host error over context-cancel-induced noise.
	if real := tracker.get(); real != nil {
		return real
	}
	if dispatchErr != nil {
		return dispatchErr
	}
	return waitErr
}

// dispatchWorkers walks hostGroup, filters by m's target group,
// acquires semaphore slots, and launches one worker per matching
// host. It returns a non-nil error only when the dispatch loop
// itself fails (e.g. pool.GetRemote error before any worker has
// run); per-host execution errors are recorded in tracker so the
// caller can prefer them over context.Canceled.
func dispatchWorkers(
	ctx context.Context,
	g *errgroup.Group,
	sem *semaphore.Weighted,
	pool *remote.HostPool,
	m module.Module,
	hostGroup map[string][]remote.Host,
	results chan<- execResult,
	tracker *firstErrorTracker,
) error {
	for group, hosts := range hostGroup {
		if m.Group() != constants.GroupAll && m.Group() != group {
			logs.Logger().Debug().
				Str("module", m.Name()).
				Str("module_group", m.Group()).
				Str("group_name", group).
				Msg("skip host group: module group does not match")
			continue
		}

		logs.Logger().Debug().
			Str("module", m.Name()).
			Str("group", group).
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
					// Context canceled by an earlier worker failure
					// or by pool.GetRemote failing in this iteration.
					// Bail out — the recorded per-host error (if any)
					// will surface via tracker.get() in ExecuteOnHosts.
					return nil
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

			g.Go(func() error {
				defer sem.Release(1)

				logs.Logger().Debug().
					Str("module", m.Name()).
					Str("host", hostStr).
					Msg("executing module on host")

				res, err := m.Execute(r)
				if err != nil {
					tracker.record(hostStr, err)
					logs.Logger().Error().
						Err(err).
						Str("module", m.Name()).
						Str("host", hostStr).
						Msg("module execution failed")
					return err
				}

				results <- execResult{h: h, res: res}
				return nil
			})
		}
	}
	return nil
}

func totalHostCount(hostGroup map[string][]remote.Host) int {
	n := 0
	for _, hosts := range hostGroup {
		n += len(hosts)
	}
	return n
}

func (e Executor) ExecuteModules(
	outputPrefix string,
	hostGroup map[string][]remote.Host,
	modules []module.Module,
) error {
	pool := remote.NewHostPool()

	for _, m := range modules {
		e.optPrinter.PrintTaskHeader(m.Name())

		if err := e.ExecuteOnHosts(outputPrefix, hostGroup, pool, m); err != nil {
			return err
		}
		e.optPrinter.PrintTaskSeparator()
	}
	return nil
}
