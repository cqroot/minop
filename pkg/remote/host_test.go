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

package remote_test

import (
	"testing"

	"github.com/cqroot/minop/pkg/remote"
	"github.com/stretchr/testify/require"
)

func TestHostsFromHostList(t *testing.T) {
	hosts, err := remote.HostsFromHostList("./testdata/host.list")
	require.NoError(t, err)
	require.Equal(t, []remote.Host{
		{"127.0.0.1", 8888, "root", "password"},
		{"127.0.0.1", 8888, "ops", "password"},
		{"127.0.0.1", 8888, "root", "pass:@:@word"},
		{"127.0.0.1", 22, "root", "password"},
	}, hosts)
}
