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

package action

import (
	"github.com/cqroot/gtypes/orderedmap"
	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
)

type Action interface {
	Validate(actCtx map[string]string) error
	Execute(h host.Host, logger *log.Logger) (*orderedmap.OrderedMap[string, string], error)
}

type ActionWrapper struct {
	Action
	role string
}

func New(role string, act Action) *ActionWrapper {
	return &ActionWrapper{
		role:   role,
		Action: act,
	}
}

func (aw *ActionWrapper) SetRole(role string) {
	aw.role = role
}

func (aw ActionWrapper) Role() string {
	return aw.role
}
