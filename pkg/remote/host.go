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
	// Username is the login username for SSH authentication.
	Username string
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
	ErrInvalidAddress       = errors.New("invalid address")
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

	h := Host{Username: user, Password: password}
	if hostPort[0] == '[' {
		h.Address, h.Port, err = parseBracketedHost(hostPort)
	} else {
		h.Address, h.Port, err = parsePlainHost(hostPort)
	}
	if err != nil {
		return Host{}, err
	}

	if h.Port < 1 || h.Port > 65535 {
		return Host{}, fmt.Errorf("%w: %d", ErrPortOutOfRange, h.Port)
	}

	return h, nil
}

// HostStr formats h as "user@addr:port" with the given prefix
// prepended (use "" for no prefix). The shape is shared by the
// executor's per-host log lines and cmd/host's tree rendering, so
// keeping it in one place ensures they stay consistent.
func HostStr(h Host, prefix string) string {
	return prefix + fmt.Sprintf("%s@%s:%d", h.Username, h.Address, h.Port)
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

// parseBracketedHost parses a "[host]" or "[host]:port" address where
// the bracketed portion must be a valid IPv4 or IPv6 literal. The
// returned address retains its surrounding brackets so it is usable
// as-is in SSH "user@[host]:port" syntax.
func parseBracketedHost(s string) (address string, port int, err error) {
	closeIdx := strings.Index(s, "]")
	if closeIdx == -1 {
		return "", 0, ErrMissingIPv6Bracket
	}
	if net.ParseIP(s[1:closeIdx]) == nil {
		return "", 0, ErrInvalidAddress
	}
	hostWithBrackets := s[:closeIdx+1]
	remaining := s[closeIdx+1:]

	if remaining != "" && remaining[0] != ':' {
		return "", 0, fmt.Errorf("%w: %s", ErrUnexpectedIPv6Suffix, remaining)
	}
	portStr := strings.TrimPrefix(remaining, ":")

	port, err = parsePort(portStr)
	if err != nil {
		return "", 0, err
	}
	return hostWithBrackets, port, nil
}

// parsePlainHost parses a "host:port" address. A bare "host" without a
// colon is accepted and port defaults to defaultPort; all other split
// failures are reported as ErrInvalidPort.
func parsePlainHost(s string) (address string, port int, err error) {
	if !strings.Contains(s, ":") {
		return s, defaultPort, nil
	}

	host, portStr, err := net.SplitHostPort(s)
	if err != nil {
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
