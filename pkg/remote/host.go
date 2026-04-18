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
	"os"
	"strings"

	"github.com/cqroot/gutils/strutils"
	"github.com/cqroot/minop/pkg/logs"
	"github.com/stretchr/testify/assert/yaml"
)

type Host struct {
	User     string
	Password string
	Address  string
	Port     int
}

var (
	ErrEmptyUsername = errors.New("empty username")
	ErrEmptyPassword = errors.New("empty password")
	ErrEmptyAddress  = errors.New("empty hostname")
	ErrInvalidPort   = errors.New("invalid port")
)

func ParseHostLine(line string) (Host, error) {
	h := Host{}
	s := line

	userDelimiter := strings.IndexByte(s, ':')
	if userDelimiter == -1 {
		return h, ErrEmptyUsername
	}
	h.User = s[:userDelimiter]
	s = s[userDelimiter+1:]

	passwordDelimiter := strings.LastIndexByte(s, '@')
	if passwordDelimiter == -1 {
		return h, ErrEmptyPassword
	}
	h.Password = s[:passwordDelimiter]
	s = s[passwordDelimiter+1:]

	if len(s) == 0 {
		return h, ErrEmptyAddress
	}

	if s[0] == '[' {
		closeIdx := strings.IndexByte(s, ']')
		if closeIdx == -1 {
			return h, fmt.Errorf("missing closing bracket for IPv6 address")
		}
		h.Address = s[:closeIdx+1]
		remaining := s[closeIdx+1:]

		if remaining == "" {
			h.Port = 22
		} else if remaining[0] == ':' {
			portStr := remaining[1:]
			if portStr == "" {
				return h, ErrInvalidPort
			}
			if !strutils.IsInteger64(portStr) {
				return h, ErrInvalidPort
			}
			h.Port = int(strutils.ToInteger64(portStr))
		} else {
			return h, fmt.Errorf("unexpected characters after IPv6 address: %s", remaining)
		}
	} else {
		hostnameDelimiter := strings.IndexByte(s, ':')
		if hostnameDelimiter == -1 {
			if len(s) != 0 {
				h.Address = s
				s = ""
			} else {
				return h, ErrEmptyAddress
			}
		} else {
			h.Address = s[:hostnameDelimiter]
			s = s[hostnameDelimiter+1:]
		}

		if s == "" {
			h.Port = 22
		} else if !strutils.IsInteger64(s) {
			return h, ErrInvalidPort
		} else {
			h.Port = int(strutils.ToInteger64(s))
		}
	}

	if h.Port < 1 || h.Port > 65535 {
		return h, fmt.Errorf("port %d not in 1-65535 range", h.Port)
	}

	return h, nil
}

func ParseHostsFile(filename string) (map[string][]Host, error) {
	logs.Logger().Debug().Str("filename", filename).Msg("Parsing hosts file")
	content, err := os.ReadFile(filename)
	if err != nil {
		logs.Logger().Err(err).Msg("")
		return nil, err
	}

	yamlContent := make(map[string][]string)
	err = yaml.Unmarshal(content, &yamlContent)
	if err != nil {
		logs.Logger().Err(err).Msg("Failed to parse hosts file as YAML")
		return nil, fmt.Errorf("failed to parse hosts file as YAML: %w", err)
	}

	hostGroup := make(map[string][]Host)
	for role, lines := range yamlContent {
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			h, err := ParseHostLine(line)
			if err != nil {
				return nil, fmt.Errorf("parse host line for role %q: %w", role, err)
			}
			hostGroup[role] = append(hostGroup[role], h)
		}
	}

	return hostGroup, nil
}
