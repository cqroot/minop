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
	// User is the login username for SSH authentication.
	User string
	// Password is the login password for SSH authentication.
	Password string
	// Address is the hostname or IP of the remote server. For IPv6, it
	// keeps the surrounding brackets (e.g. "[2001:db8::1]") so it can
	// be used directly with SSH address syntax.
	Address string
	// Port is the TCP port of the remote server (1-65535).
	Port int
}

// defaultPort is the SSH port used when the input does not specify one.
const defaultPort = 22

// Host parsing errors returned by ParseHostLine.
var (
	ErrEmptyUsername        = errors.New("empty username")
	ErrEmptyPassword        = errors.New("empty password")
	ErrEmptyAddress         = errors.New("empty hostname")
	ErrInvalidPort          = errors.New("invalid port")
	ErrMissingIPv6Bracket   = errors.New("missing closing bracket for IPv6 address")
	ErrUnexpectedIPv6Suffix = errors.New("unexpected characters after IPv6 address")
	ErrPortOutOfRange       = errors.New("port out of range")
)

// ParseHostLine parses a single host line in the form
// "user:password@host:port" and returns the resulting Host. Supported
// shapes for the address part are:
//
//   - "host"             — defaults port to defaultPort
//   - "host:port"        — host with explicit port
//   - "[ipv6]"           — IPv6 without port, defaults port
//   - "[ipv6]:port"      — IPv6 with explicit port
//
// The rightmost '@' is treated as the user/host separator so that
// passwords may contain '@'. Returns one of the sentinel errors
// declared in this package on parse failure.
func ParseHostLine(line string) (Host, error) {
	user, password, hostPort, err := parseUserInfo(line)
	if err != nil {
		return Host{}, err
	}

	if hostPort == "" {
		return Host{}, ErrEmptyAddress
	}

	h := Host{User: user, Password: password}
	if hostPort[0] == '[' {
		h.Address, h.Port, err = parseIPv6(hostPort)
	} else {
		h.Address, h.Port, err = parseHostPort(hostPort)
	}
	if err != nil {
		return Host{}, err
	}

	if h.Port < 1 || h.Port > 65535 {
		return Host{}, fmt.Errorf("%w: %d", ErrPortOutOfRange, h.Port)
	}

	return h, nil
}

// parseUserInfo splits "user:password@hostPort" into its components.
// The rightmost '@' is used as the separator, so passwords may
// themselves contain '@'.
func parseUserInfo(s string) (user, password, hostPort string, err error) {
	atIdx := strings.LastIndexByte(s, '@')
	if atIdx == -1 {
		return "", "", "", ErrEmptyPassword
	}

	userPart := s[:atIdx]
	hostPort = s[atIdx+1:]

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

	return user, password, hostPort, nil
}

// parseIPv6 parses a "[ipv6]" or "[ipv6]:port" address. The returned
// address retains its surrounding brackets so it is usable as-is in
// SSH "user@[ipv6]:port" syntax.
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

// parseHostPort parses a "host:port" address. A bare "host" without a
// colon is accepted and port defaults to defaultPort; all other split
// failures are reported as ErrInvalidPort.
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

// parsePort returns defaultPort when portStr is empty, otherwise it
// parses portStr as an integer. A non-integer portStr yields
// ErrInvalidPort.
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
