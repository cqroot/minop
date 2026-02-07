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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/operation"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"golang.org/x/term"
)

type Executor struct {
	optVerboseLevel int
	optMaxProcs     int
	outputPrefix    string
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

func (e Executor) printValue(key string, val string) {
	if val == "" {
		return
	}

	prefix := fmt.Sprintf("%s    ", e.outputPrefix)
	if e.optVerboseLevel == 0 && (strings.IndexByte(val, '\n') == -1 || strings.IndexByte(val, '\n') == len(val)-1) {
		fmt.Printf("%s%s %s\n", prefix, color.CyanString("%s:", key), strings.ReplaceAll(val, "\n", ""))
	} else {
		color.Cyan("%s%s:\n", prefix, key)
		scanner := bufio.NewScanner(strings.NewReader(val))
		for scanner.Scan() {
			fmt.Printf("%s    %s\n", prefix, scanner.Text())
		}
	}
}

type execResult struct {
	h   remote.Host
	res *gtypes.OrderedMap[string, string]
}

func (e Executor) ExecuteOperation(hostGroup map[string][]remote.Host, pool *remote.HostPool, op operation.Operation) error {
	chanExecResults := make(chan execResult)

	printDone := make(chan struct{})
	go func() {
		defer close(printDone)
		for res := range chanExecResults {
			hostStr := fmt.Sprintf("%s%s@%s:%d", e.outputPrefix, res.h.User, res.h.Address, res.h.Port)
			fmt.Printf("%s  %s\n", color.HiCyanString(hostStr),
				color.HiBlackString(time.Now().Format("[2006-01-02 15:04:05]")))

			if res.res != nil {
				_ = res.res.ForEach(func(key, val string) error {
					e.printValue(key, val)
					return nil
				})
			}
		}
	}()

	sem := semaphore.NewWeighted(int64(e.optMaxProcs))
	g, ctx := errgroup.WithContext(context.Background())

	for role, hosts := range hostGroup {
		if op.Role() != "all" && op.Role() != role {
			continue
		}

		for _, h := range hosts {
			if err := sem.Acquire(ctx, 1); err != nil {
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

func (e Executor) ExecuteOperations(ops []operation.Operation) error {
	hostGroup, err := remote.ParseHostsFile(filepath.Join(".", constants.HostFileName))
	if err != nil {
		return err
	}

	pool := remote.NewHostPool()
	e.outputPrefix = "    "

	for _, op := range ops {
		termWidth := 500
		if term.IsTerminal(int(os.Stdout.Fd())) {
			width, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err == nil {
				termWidth = width
			}
		}

		fmt.Printf("%s %s %s\n",
			color.HiCyanString(op.Name()),
			color.HiBlackString(strings.Repeat("•", termWidth-len(op.Name())-2-19)),
			color.HiBlackString(time.Now().Format("2006-01-02 15:04:05")),
		)

		err := e.ExecuteOperation(hostGroup, pool, op)
		if err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}
