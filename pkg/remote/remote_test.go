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

// TestRemoteClose_NilWhenNothingToClose verifies that calling Close
// on a *Remote built via NewForTesting (no real SSH/SFTP clients)
// returns nil without panicking. This guards the early-exit path
// inside Close after the refactor to errors.Join.
func TestRemoteClose_NilWhenNothingToClose(t *testing.T) {
	h, err := remote.ParseHostLine("user:pw@127.0.0.1:22")
	require.NoError(t, err)

	r := remote.NewForTesting(h)
	require.NoError(t, r.Close())
}
