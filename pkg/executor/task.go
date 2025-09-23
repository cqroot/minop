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
	"bufio"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cqroot/gtypes/orderedmap"
	"github.com/cqroot/minop/pkg/action"
	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type ActionExecutor struct {
	logger          *log.Logger
	optVerboseLevel int
	optMaxProcs     int
}

func New(logger *log.Logger, opts ...Option) *ActionExecutor {
	e := ActionExecutor{
		logger:          logger,
		optVerboseLevel: 0,
		optMaxProcs:     1,
	}

	for _, opt := range opts {
		opt(&e)
	}

	return &e
}

func (e ActionExecutor) printValue(key string, val string, prefix string) {
	if e.optVerboseLevel == 0 && (strings.IndexByte(val, '\n') == -1 || strings.IndexByte(val, '\n') == len(val)-1) {
		if val != "" {
			fmt.Printf("%s%s %s\n", prefix, color.CyanString("%s:", key), strings.ReplaceAll(val, "\n", ""))
		}
	} else {
		color.Cyan("%s%s:\n", prefix, key)
		scanner := bufio.NewScanner(strings.NewReader(val))
		for scanner.Scan() {
			fmt.Printf("%s    %s\n", prefix, scanner.Text())
		}
	}
}

type execResult struct {
	h   host.Host
	res *orderedmap.OrderedMap[string, string]
}

func (e ActionExecutor) PrintActionResult(hostGroup map[string][]host.Host, rgs *map[host.Host]*remote.Remote, act action.ActionWrapper, prefix string) error {
	chanExecResults := make(chan execResult)

	printDone := make(chan struct{})
	go func() {
		defer close(printDone)
		for res := range chanExecResults {
			hostStr := fmt.Sprintf("%s%s@%s:%d", prefix, res.h.User, res.h.Address, res.h.Port)
			fmt.Printf("%s  %s\n", color.HiCyanString(hostStr),
				color.HiBlackString(time.Now().Format("[2006-01-02 15:04:05]")))

			if res.res != nil {
				_ = res.res.ForEach(func(key, val string) error {
					e.printValue(key, val, fmt.Sprintf("%s    ", prefix))
					return nil
				})
			}
		}
	}()

	sem := semaphore.NewWeighted(int64(e.optMaxProcs))
	g, ctx := errgroup.WithContext(context.Background())

	for role, hosts := range hostGroup {
		if act.Role() != "all" && act.Role() != role {
			continue
		}

		for _, h := range hosts {
			if err := sem.Acquire(ctx, 1); err != nil {
				continue
			}

			r, ok := (*rgs)[h]
			if !ok {
				newR, err := remote.New(h, e.logger)
				if err != nil {
					return err
				}
				(*rgs)[h] = newR
				r = newR
			}

			currHost := h
			g.Go(func() error {
				defer sem.Release(1)

				res, err := act.Execute(r, e.logger)
				if err != nil {
					return err
				}

				chanExecResults <- execResult{
					h:   currHost,
					res: res,
				}
				return nil
			})
		}
	}

	var firstErr error
	go func() {
		err := g.Wait()
		if err != nil {
			firstErr = err
		}
		close(chanExecResults)
	}()

	<-printDone
	return firstErr
}
