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

	"github.com/cqroot/minop/pkg/utils"

	"gopkg.in/yaml.v3"
)

type Host struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	ErrEmptyLine     = errors.New("empty line")
	ErrEmptyUsername = errors.New("empty username")
	ErrEmptyPassword = errors.New("empty password")
	ErrEmptyHostname = errors.New("empty hostname")
)

func HostFromLine(line string) (Host, error) {
	h := Host{}
	if len(line) == 0 {
		return h, ErrEmptyLine
	}

	s := line

	userDelimiter := strings.IndexByte(s, ':')
	if userDelimiter == -1 {
		return h, fmt.Errorf("%w: %s", ErrEmptyUsername, line)
	} else {
		h.Username = s[:userDelimiter]
		s = s[userDelimiter+1:]
	}

	passwordDelimiter := strings.LastIndexByte(s, '@')
	if passwordDelimiter == -1 {
		return h, fmt.Errorf("%w: %s", ErrEmptyPassword, line)
	} else {
		h.Password = s[:passwordDelimiter]
		s = s[passwordDelimiter+1:]
	}

	hostnameDelimiter := strings.IndexByte(s, ':')
	if hostnameDelimiter == -1 {
		if len(s) != 0 {
			h.Hostname = s
			s = ""
		} else {
			return h, fmt.Errorf("%w: %s", ErrEmptyHostname, line)
		}
	} else {
		h.Hostname = s[:hostnameDelimiter]
		s = s[hostnameDelimiter+1:]
	}

	if !utils.StrIsInteger(s) {
		h.Port = 22
	} else {
		h.Port = int(utils.StrToInteger(s))
	}

	return h, nil
}

func HostsFromHostList(filename string) ([]Host, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	hosts := make([]Host, 0)
	for fileScanner.Scan() {
		h, err := HostFromLine(fileScanner.Text())
		if errors.Is(err, ErrEmptyLine) {
			continue
		}

		if err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}
	return hosts, nil
}

func HostsFromYaml(filename string) ([]Host, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	hosts := make([]Host, 0)
	err = yaml.Unmarshal(content, &hosts)
	if err != nil {
		return nil, err
	}

	for i := range hosts {
		if hosts[i].Port == 0 {
			hosts[i].Port = 22
		}
		if hosts[i].Username == "" {
			hosts[i].Username = "root"
		}
	}
	return hosts, nil
}
