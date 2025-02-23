package manager

import (
	"errors"
	"fmt"
	"github.com/cqroot/minop/pkg/module"
	"github.com/cqroot/minop/pkg/module/command"
	"github.com/cqroot/minop/pkg/remote"
)

var (
	ErrModuleParse = errors.New("module parse error")
)

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
	for _, m := range mgr.modules {
		_, err := m.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
