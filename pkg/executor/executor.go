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
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/operation"
	"github.com/cqroot/minop/pkg/remote"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"golang.org/x/term"
)

var (
	labelStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	taskStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	dimStyle       = lipgloss.NewStyle().Faint(true)
	hostStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	timestampStyle = lipgloss.NewStyle().Faint(true)
)

// Executor orchestrates remote operations across multiple hosts.
type Executor struct {
	optVerboseLevel int
	optMaxProcs     int
	outputPrefix    string
}

// New creates a new Executor with the given options.
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

// printValue outputs a key-value pair with the configured prefix.
// Single-line values are printed inline; multi-line values are printed block-style.
func (e Executor) printValue(key string, val string) {
	if val == "" {
		return
	}

	prefix := fmt.Sprintf("%s    ", e.outputPrefix)
	if e.optVerboseLevel == 0 && (strings.IndexByte(val, '\n') == -1 || strings.IndexByte(val, '\n') == len(val)-1) {
		fmt.Printf("%s%s %s\n", prefix, labelStyle.Render(fmt.Sprintf("%s:", key)), strings.ReplaceAll(val, "\n", ""))
	} else {
		fmt.Printf("%s%s:\n", prefix, labelStyle.Render(key))
		scanner := bufio.NewScanner(strings.NewReader(val))
		for scanner.Scan() {
			fmt.Printf("%s    %s\n", prefix, scanner.Text())
		}
	}
}

// execResult holds the result of a remote operation execution.
type execResult struct {
	h   remote.Host
	res *gtypes.OrderedMap[string, string]
}

// ExecuteOperation runs a single operation on all matching hosts in the group.
// It respects the operation's Role field: if Role is "all", it runs on all hosts;
// otherwise, it runs only on hosts in the specified role group.
func (e Executor) ExecuteOperation(hostGroup map[string][]remote.Host, pool *remote.HostPool, op operation.Operation) error {
	execResultsChan := make(chan execResult)

	printDone := make(chan struct{})
	go func() {
		defer close(printDone)
		for res := range execResultsChan {
			hostStr := fmt.Sprintf("%s%s@%s:%d", e.outputPrefix, res.h.User, res.h.Address, res.h.Port)
			fmt.Printf("%s  %s\n", hostStyle.Render(hostStr),
				timestampStyle.Render(time.Now().Format("[2006-01-02 15:04:05]")))

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

				execResultsChan <- execResult{
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
		close(execResultsChan)
	}()

	<-printDone
	return <-errCh
}

// ExecuteOperations runs a sequence of operations on the host group.
// Each operation is executed on all hosts that match the operation's Role.
func (e Executor) ExecuteOperations(hostGroup map[string][]remote.Host, ops []operation.Operation) error {
	pool := remote.NewHostPool()
	e.outputPrefix = "    "

	termWidth := 500
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			termWidth = w
		}
	}

	for _, op := range ops {
		delim := ""
		delimLen := termWidth - len(op.Name()) - 2 - 19
		if delimLen > 0 {
			delim = strings.Repeat("•", delimLen)
		}
		fmt.Printf("%s %s %s\n",
			taskStyle.Render(op.Name()),
			dimStyle.Render(delim),
			dimStyle.Render(time.Now().Format("2006-01-02 15:04:05")),
		)

		err := e.ExecuteOperation(hostGroup, pool, op)
		if err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}
