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

const defaultPort = 22

// Host parsing errors
var (
	ErrEmptyUsername        = errors.New("empty username")
	ErrEmptyPassword        = errors.New("empty password")
	ErrEmptyAddress         = errors.New("empty hostname")
	ErrInvalidPort          = errors.New("invalid port")
	ErrMissingIPv6Bracket   = errors.New("missing closing bracket for IPv6 address")
	ErrUnexpectedIPv6Suffix = errors.New("unexpected characters after IPv6 address")
	ErrPortOutOfRange       = errors.New("port out of range")
)

func ParseHostLine(line string) (Host, error) {
	user, password, rest, err := parseUserInfo(line)
	if err != nil {
		return Host{}, err
	}

	if rest == "" {
		return Host{}, ErrEmptyAddress
	}

	h := Host{User: user, Password: password}
	if rest[0] == '[' {
		h.Address, h.Port, err = parseIPv6(rest)
	} else {
		h.Address, h.Port, err = parseHostPort(rest)
	}
	if err != nil {
		return Host{}, err
	}

	if h.Port < 1 || h.Port > 65535 {
		return Host{}, fmt.Errorf("%w: %d", ErrPortOutOfRange, h.Port)
	}

	return h, nil
}

// parseUserInfo extracts user and password from "user:password@...".
func parseUserInfo(s string) (user, password, rest string, err error) {
	atIdx := strings.LastIndexByte(s, '@')
	if atIdx == -1 {
		return "", "", "", ErrEmptyPassword
	}

	userPart := s[:atIdx]
	rest = s[atIdx+1:]

	colonIdx := strings.IndexByte(userPart, ':')
	if colonIdx == -1 {
		return "", "", "", ErrEmptyUsername
	}
	user = userPart[:colonIdx]
	password = userPart[colonIdx+1:]
	if user == "" {
		return "", "", "", ErrEmptyUsername
	}
	if password == "" {
		return "", "", "", ErrEmptyPassword
	}

	return user, password, rest, nil
}

// parseIPv6 parses a "[ipv6]:port" address.
func parseIPv6(s string) (address string, port int, err error) {
	closeIdx := strings.Index(s, "]")
	if closeIdx == -1 {
		return "", 0, ErrMissingIPv6Bracket
	}
	hostWithBrackets := s[:closeIdx+1]
	remaining := s[closeIdx+1:]

	var portStr string
	if remaining == "" {
		portStr = ""
	} else if remaining[0] == ':' {
		portStr = remaining[1:]
	} else {
		return "", 0, fmt.Errorf("%w: %s", ErrUnexpectedIPv6Suffix, remaining)
	}

	port, err = parsePort(portStr)
	if err != nil {
		return "", 0, err
	}
	return hostWithBrackets, port, nil
}

// parseHostPort parses a "host:port" address.
func parseHostPort(s string) (address string, port int, err error) {
	host, portStr, err := net.SplitHostPort(s)
	if err != nil {
		var addrErr *net.AddrError
		if errors.As(err, &addrErr) && addrErr.Err == "missing port in address" {
			return s, defaultPort, nil
		}
		return "", 0, ErrInvalidPort
	}

	if host == "" {
		return "", 0, ErrEmptyAddress
	}

	port, err = parsePort(portStr)
	if err != nil {
		return "", 0, err
	}
	return host, port, nil
}

// parsePort returns defaultPort if portStr is empty, otherwise parses it as an integer.
func parsePort(portStr string) (int, error) {
	if portStr == "" {
		return defaultPort, nil
	}
	p, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, ErrInvalidPort
	}
	return p, nil
}
