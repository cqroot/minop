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

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Host represents a remote server connection with authentication details.
type Host struct {
	User     string
	Password string
	Address  string
	Port     int
}

// Host parsing errors
var (
	ErrEmptyUsername      = errors.New("empty username")
	ErrEmptyPassword      = errors.New("empty password")
	ErrEmptyAddress       = errors.New("empty hostname")
	ErrInvalidPort        = errors.New("invalid port")
	ErrMissingIPv6Bracket = errors.New("missing closing bracket for IPv6 address")
)

func ParseHostLine(line string) (Host, error) {
	h := Host{}
	s := line

	atIdx := strings.LastIndexByte(s, '@')
	if atIdx == -1 {
		return Host{}, ErrEmptyPassword
	}

	userPart := s[:atIdx]
	s = s[atIdx+1:]

	colonIdx := strings.IndexByte(userPart, ':')
	if colonIdx == -1 {
		return Host{}, ErrEmptyUsername
	}
	h.User = userPart[:colonIdx]
	h.Password = userPart[colonIdx+1:]
	if h.User == "" {
		return Host{}, ErrEmptyUsername
	}
	if h.Password == "" {
		return Host{}, ErrEmptyPassword
	}

	if s == "" {
		return Host{}, ErrEmptyAddress
	}

	if s[0] == '[' {
		closeIdx := strings.Index(s, "]")
		if closeIdx == -1 {
			return Host{}, ErrMissingIPv6Bracket
		}
		hostWithBrackets := s[:closeIdx+1]
		remaining := s[closeIdx+1:]

		var portStr string
		if remaining == "" {
			portStr = ""
		} else if remaining[0] == ':' {
			portStr = remaining[1:]
		} else {
			return Host{}, fmt.Errorf("unexpected characters after IPv6 address: %s", remaining)
		}

		h.Address = hostWithBrackets
		if portStr == "" {
			h.Port = 22
		} else {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return Host{}, ErrInvalidPort
			}
			h.Port = port
		}
	} else {
		host, portStr, err := net.SplitHostPort(s)
		if err != nil {
			var addrErr *net.AddrError
			if errors.As(err, &addrErr) && addrErr.Err == "missing port in address" {
				h.Address = s
				h.Port = 22
			} else {
				return Host{}, ErrInvalidPort
			}
		} else {
			h.Address = host
			if h.Address == "" {
				return Host{}, ErrEmptyAddress
			}
			if portStr == "" {
				h.Port = 22
			} else {
				port, err := strconv.Atoi(portStr)
				if err != nil {
					return Host{}, ErrInvalidPort
				}
				h.Port = port
			}
		}
	}

	if h.Port < 1 || h.Port > 65535 {
		return Host{}, fmt.Errorf("port %d not in 1-65535 range", h.Port)
	}

	return h, nil
}
