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

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/operation"
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
		optMaxProcs:     1,
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

func (e Executor) ExecuteOnHosts(outputPrefix string, hostGroup map[string][]remote.Host, pool *remote.HostPool, op operation.Operation) error {
	if localOp, ok := op.(*operation.OpLocal); ok {
		localOp.SetPrefix(outputPrefix)
		res, err := localOp.Execute(nil)
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

	for role, hosts := range hostGroup {
		if op.Role() != constants.RoleAll && op.Role() != role {
			continue
		}

		for _, h := range hosts {
			if err := sem.Acquire(ctx, 1); err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				continue
			}

			r, err := pool.GetRemote(h)
			if err != nil {
				return err
			}

			currHost := h
			g.Go(func() error {
				defer sem.Release(1)

				res, err := op.Execute(r)
				if err != nil {
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
	return <-errCh
}

func (e Executor) ExecuteOperations(outputPrefix string, hostGroup map[string][]remote.Host, ops []operation.Operation) error {
	termWidth := getTerminalWidth()
	pool := remote.NewHostPool()

	for _, op := range ops {
		printTaskHeader(op.Name(), termWidth)

		if err := e.ExecuteOnHosts(outputPrefix, hostGroup, pool, op); err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}
