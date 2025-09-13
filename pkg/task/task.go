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

package task

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cqroot/gtypes/orderedmap"
	"github.com/cqroot/minop/pkg/action"
	"github.com/cqroot/minop/pkg/action/command"
	"github.com/cqroot/minop/pkg/action/file"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
	"github.com/cqroot/minop/pkg/utils/maputils"
	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"gopkg.in/yaml.v3"
)

type Task struct {
	actions         []action.ActionWrapper
	logger          *log.Logger
	optVerboseLevel int
	optMaxProcs     int
}

func New(filename string, logger *log.Logger, opts ...Option) (*Task, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var actCtxs []map[string]string
	err = yaml.Unmarshal(content, &actCtxs)
	if err != nil {
		return nil, err
	}

	acts := make([]action.ActionWrapper, len(actCtxs))

	for i, actCtx := range actCtxs {
		actName, err := maputils.GetString(actCtx, "action")
		if err != nil {
			return nil, err
		}

		role := maputils.GetStringOrDefault(actCtx, "role", "all")
		var act action.Action
		switch actName {
		case "command":
			act, err = command.New(actCtx)
			if err != nil {
				return nil, err
			}
		case "file":
			act, err = file.New(actCtx)
			if err != nil {
				return nil, err
			}
		}
		acts[i] = *action.New(maputils.GetStringOrDefault(actCtx, "name", actName), role, act)
	}

	t := Task{
		actions:         acts,
		logger:          logger,
		optVerboseLevel: 0,
		optMaxProcs:     1,
	}
	for _, opt := range opts {
		opt(&t)
	}
	return &t, nil
}

func (t Task) printValue(key string, val string, prefix string) {
	if t.optVerboseLevel == 1 && (strings.IndexByte(val, '\n') == -1 || strings.IndexByte(val, '\n') == len(val)-1) {
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

func (t Task) PrintActionResult(hostGroup map[string][]host.Host, act action.ActionWrapper, prefix string) error {
	chanExecResults := make(chan execResult)

	printDone := make(chan struct{})
	go func() {
		defer close(printDone)
		for res := range chanExecResults {
			color.HiCyan("%s%s@%s:%d", prefix, res.h.User, res.h.Address, res.h.Port)

			if t.optVerboseLevel == 0 {
				continue
			}

			if res.res != nil {
				_ = res.res.ForEach(func(key, val string) error {
					t.printValue(key, val, fmt.Sprintf("%s    ", prefix))
					return nil
				})
			}
		}
	}()

	sem := semaphore.NewWeighted(int64(t.optMaxProcs))
	g, ctx := errgroup.WithContext(context.Background())

	for role, hosts := range hostGroup {
		if act.Role() != "all" && act.Role() != role {
			continue
		}

		for _, h := range hosts {
			if err := sem.Acquire(ctx, 1); err != nil {
				continue
			}

			currHost := h
			g.Go(func() error {
				defer sem.Release(1)

				res, err := act.Execute(currHost, t.logger)
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

func (t Task) Execute() error {
	hostGroup, err := host.Read(filepath.Join(".", constants.HostFileName))
	if err != nil {
		return err
	}

	for _, act := range t.actions {
		color.HiCyan("%s:\n", act.Name())
		err := t.PrintActionResult(hostGroup, act, "    ")
		if err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}
