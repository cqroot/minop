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

package manager

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cqroot/minop/pkg/module"
	"github.com/cqroot/minop/pkg/module/command"
	"github.com/cqroot/minop/pkg/remote"
)

var ErrModuleParse = errors.New("module parse error")

type Manager struct {
	moduleConfig *ModuleConfig
	rem          *remote.Remote
	modules      []module.Module
}

func New(host remote.Host, moduleArgs *ModuleConfig) (*Manager, error) {
	mgr := Manager{
		moduleConfig: moduleArgs,
	}

	r, err := remote.New(host)
	if err != nil {
		return nil, err
	}
	mgr.rem = r

	for _, moduleArg := range mgr.moduleConfig.ModuleArgs {
		name, ok := moduleArg["name"]
		if !ok {
			return nil, fmt.Errorf("%w: %+v", ErrModuleParse, moduleArg)
		}

		switch name {
		case "command":
			m, err := command.New(mgr.rem, moduleArg)
			if err != nil {
				return nil, err
			}
			mgr.modules = append(mgr.modules, m)
		}
	}

	return &mgr, nil
}

func (mgr *Manager) Close() {
	_ = mgr.rem.Close()
}

func (mgr *Manager) Run() error {
	resultsCh := make(chan string)
	var wgProcess sync.WaitGroup
	var wgPrint sync.WaitGroup

	wgPrint.Add(1)
	go func() {
		defer wgPrint.Done()
		for result := range resultsCh {
			fmt.Println(result)
		}
	}()

	for _, m := range mgr.modules {
		wgProcess.Add(1)
		go func() {
			defer wgProcess.Done()

			_ = m.Run(resultsCh)
		}()
		wgProcess.Wait()
	}
	close(resultsCh)

	wgProcess.Wait()
	return nil
}
