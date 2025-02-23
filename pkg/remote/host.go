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
	"gopkg.in/yaml.v3"
	"os"
)

type Host struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
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

	for i, _ := range hosts {
		if hosts[i].Port == 0 {
			hosts[i].Port = 22
		}
		if hosts[i].Username == "" {
			hosts[i].Username = "root"
		}
	}
	return hosts, nil
}
