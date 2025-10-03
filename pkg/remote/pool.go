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

package remote

import "github.com/rs/zerolog"

type HostPool struct {
	hosts  map[Host]*Remote
	logger zerolog.Logger
}

func NewHostPool(logger zerolog.Logger) *HostPool {
	return &HostPool{
		hosts:  make(map[Host]*Remote),
		logger: logger,
	}
}

func (p *HostPool) GetRemote(host Host) (*Remote, error) {
	r, ok := p.hosts[host]
	if !ok {
		newR, err := New(host, p.logger)
		if err != nil {
			return nil, err
		}
		p.hosts[host] = newR
		r = newR
	}
	return r, nil
}
