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
	"fmt"
	"testing"

	"github.com/cqroot/minop/pkg/remote"
	"github.com/stretchr/testify/require"
)

type ParseHostLineTestCase struct {
	name     string
	line     string
	expected remote.Host
	err      error
}

func TestParseHostLine(t *testing.T) {
	tests := []ParseHostLineTestCase{
		{
			name: "valid host line",
			line: "user:password@hostname:22",
			expected: remote.Host{
				User:     "user",
				Password: "password",
				Address:  "hostname",
				Port:     22,
			},
			err: nil,
		},
		{
			name:     "empty username",
			line:     ":password@hostname:22",
			expected: remote.Host{},
			err:      remote.ErrEmptyUsername,
		},
		{
			name:     "empty password",
			line:     "user:@hostname:22",
			expected: remote.Host{},
			err:      remote.ErrEmptyPassword,
		},
		{
			name:     "empty hostname",
			line:     "user:password@:22",
			expected: remote.Host{},
			err:      remote.ErrEmptyAddress,
		},
		{
			name:     "invalid port",
			line:     "user:password@hostname:notaport",
			expected: remote.Host{},
			err:      remote.ErrInvalidPort,
		},
		{
			name:     "missing IPv6 closing bracket",
			line:     "user:password@[2001:db8::1:22",
			expected: remote.Host{},
			err:      remote.ErrMissingIPv6Bracket,
		},
		{
			name:     "unexpected characters after IPv6 address",
			line:     "user:password@[2001:db8::1]extra:22",
			expected: remote.Host{},
			err:      fmt.Errorf("unexpected characters after IPv6 address: extra:22"),
		},
		{
			name:     "port out of range",
			line:     "user:password@hostname:70000",
			expected: remote.Host{},
			err:      fmt.Errorf("port 70000 not in 1-65535 range"),
		},
		{
			name: "IPv6 address with port",
			line: "user:password@[2001:db8::1]:22",
			expected: remote.Host{
				User:     "user",
				Password: "password",
				Address:  "[2001:db8::1]",
				Port:     22,
			},
			err: nil,
		},
		{
			name: "IPv6 address without port",
			line: "user:password@[2001:db8::1]",
			expected: remote.Host{
				User:     "user",
				Password: "password",
				Address:  "[2001:db8::1]",
				Port:     22,
			},
			err: nil,
		},
		{
			name: "hostname without port",
			line: "user:password@hostname",
			expected: remote.Host{
				User:     "user",
				Password: "password",
				Address:  "hostname",
				Port:     22,
			},
			err: nil,
		},
		{
			name: "hostname with empty port",
			line: "user:password@hostname:",
			expected: remote.Host{
				User:     "user",
				Password: "password",
				Address:  "hostname",
				Port:     22,
			},
			err: nil,
		},
		{
			name: "IPv6 address with empty port",
			line: "user:password@[2001:db8::1]:",
			expected: remote.Host{
				User:     "user",
				Password: "password",
				Address:  "[2001:db8::1]",
				Port:     22,
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := remote.ParseHostLine(tt.line)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}
