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

package cmd

import (
	"os"

	"github.com/cqroot/minop/pkg/manager"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/minop/pkg/utils"
	"github.com/spf13/cobra"
)

func CheckErr(err error) {
	if err != nil {
		utils.LogError(err.Error())
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "minop",
		Short: "MinOP is a simple remote execution and deployment tool",
		Long:  "MinOP is a simple remote execution and deployment tool",
		Run: func(cmd *cobra.Command, args []string) {
			var hosts []remote.Host
			var err error
			if utils.IsFileExists("./host.list") {
				hosts, err = remote.HostsFromHostList("./host.list")
				CheckErr(err)
			} else if utils.IsFileExists("./hosts.yaml") {
				hosts, err = remote.HostsFromYaml("./hosts.yaml")
				CheckErr(err)
			}

			moduleConfig, err := manager.ModulesFromYaml("./minop.yaml")
			if err != nil {
				utils.LogError(err.Error())
				os.Exit(1)
			}

			for _, host := range hosts {
				mgr, err := manager.New(host, moduleConfig)
				if err != nil {
					utils.LogError(err.Error())
					continue
				}
				defer mgr.Close()

				err = mgr.Run()
				if err != nil {
					utils.LogError(err.Error())
				}
			}
		},
	}

	cmd.AddCommand(NewDocCmd())
	return &cmd
}

func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}
