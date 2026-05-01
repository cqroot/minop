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

// HostPool manages a cache of Remote connections keyed by Host.
// It reuses existing connections to avoid redundant SSH/SFTP handshakes.
type HostPool struct {
	hosts map[Host]*Remote
}

// NewHostPool creates a new empty HostPool.
func NewHostPool() *HostPool {
	return &HostPool{
		hosts: make(map[Host]*Remote),
	}
}

// GetRemote returns a Remote connection for the given Host.
// If a connection already exists in the pool, it returns the cached one.
// Otherwise, it creates a new connection and caches it.
func (p *HostPool) GetRemote(host Host) (*Remote, error) {
	r, ok := p.hosts[host]
	if !ok {
		newR, err := New(host)
		if err != nil {
			return nil, err
		}
		p.hosts[host] = newR
		r = newR
	}
	return r, nil
}
