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
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/minop/pkg/theme"
	"github.com/spf13/cobra"
)

// renderHost formats a host as 'user@addr:port' with each segment
// in a distinct color (user bold-cyan, address green, port blue,
// separators dim) so the three pieces are visually separable at a
// glance.
func renderHost(h remote.Host) string {
	return theme.HostUser().Render(h.Username) +
		theme.HostSep().Render("@") +
		theme.HostAddr().Render(h.Address) +
		theme.HostSep().Render(":") +
		theme.HostPort().Render(strconv.Itoa(h.Port))
}

// RunHostCmd displays all configured hosts in a tree format.
func RunHostCmd(cmd *cobra.Command, args []string) {
	e := executor.New(
		executor.WithVerboseLevel(flagVerboseLevel),
		executor.WithMaxProcs(flagMaxProcs),
	)

	hostGroup, err := e.LoadHostsFile(flagHostsFile)
	CheckErr(err)

	printHostTree(os.Stdout, hostGroup)
}

// printHostTree renders hostGroup as an indented tree to w. Groups
// are written in the order returned by that helper — the "all" role
// first, then any remaining groups alphabetically — with a blank
// line between groups for readability. Groups other than "all" are
// rendered with an extra indent level so they visually nest under
// the top-level "all" group.
func printHostTree(w io.Writer, hostGroup map[string][]remote.Host) {
	groups := sortedGroupNames(hostGroup)

	_, _ = fmt.Fprintln(w)
	for idx, group := range groups {
		printHostGroup(w, hostGroup[group], group)
		if idx < len(groups)-1 {
			_, _ = fmt.Fprintln(w)
		}
	}
}

// sortedGroupNames returns the keys of hostGroup. The special "all"
// group is placed first; the remaining groups are sorted
// alphabetically for stable, deterministic output.
func sortedGroupNames(hostGroup map[string][]remote.Host) []string {
	names := make([]string, 0, len(hostGroup))
	for name := range hostGroup {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		if names[i] == "all" {
			return true
		}
		if names[j] == "all" {
			return false
		}
		return names[i] < names[j]
	})
	return names
}

// printHostGroup renders a single group header followed by its
// hosts as tree branches. The "all" group sits at the base indent;
// any other group is rendered with an extra indent level so they
// visually nest under "all". The last host in a group uses └──
// instead of ├── to close the branch.
func printHostGroup(w io.Writer, hosts []remote.Host, name string) {
	indent := "  "
	if name != "all" {
		indent = "    "
	}

	_, _ = fmt.Fprintf(w, "%s%s\n", indent, theme.HostHeader().Render("• "+name))
	for i, host := range hosts {
		branch := theme.TreeBranch().Render("├──")
		if i == len(hosts)-1 {
			branch = theme.TreeBranch().Render("└──")
		}
		_, _ = fmt.Fprintf(w, "%s%s %s\n",
			indent,
			branch,
			renderHost(host))
	}
}

// NewHostCmd creates the host command that lists all configured hosts.
func NewHostCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "host",
		Short: "List all hosts",
		Long:  "List all hosts.",
		Run:   RunHostCmd,
	}

	return &c
}
