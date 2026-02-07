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
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cqroot/gutils/strutils"
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
	} else {
		h.User = s[:userDelimiter]
		s = s[userDelimiter+1:]
	}

	passwordDelimiter := strings.LastIndexByte(s, '@')
	if passwordDelimiter == -1 {
		return h, ErrEmptyPassword
	} else {
		h.Password = s[:passwordDelimiter]
		s = s[passwordDelimiter+1:]
	}

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
	if h.Port < 1 || h.Port > 65535 {
		return h, fmt.Errorf("port %d not in 1-65535 range", h.Port)
	}

	return h, nil
}

func ParseHostsFile(filename string) (map[string][]Host, error) {
	hostGroup := make(map[string][]Host)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	currGroup := "default"
	lineNum := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineNum++

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		if len(trimmed) >= 3 && strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			currGroup = trimmed[1 : len(trimmed)-1]
		} else {
			h, err := ParseHostLine(trimmed)
			if err != nil {
				return nil, fmt.Errorf("parse host lineline %d (%q): %w", lineNum, line, err)
			}
			hostGroup[currGroup] = append(hostGroup[currGroup], h)
		}
	}

	return hostGroup, nil
}
