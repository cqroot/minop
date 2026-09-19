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

// Package theme centralises the lipgloss styles used across the
// command, executor, CLI, and version packages. The styles are
// grouped semantically — "this is the brand colour", "this is a
// secondary label" — rather than by colour number, so callers don't
// need to memorise ANSI codes and a future palette swap only touches
// this file.
package theme

import "github.com/charmbracelet/lipgloss"

// Semantic colour palette. Each constant maps to an ANSI 256-bit
// colour used by the styles below; the names describe the role
// rather than the hue, so a palette change is a one-file edit.
const (
	colorBrand = "5"   // magenta — brand / titles
	colorTask  = "10"  // green  — task names, help keys
	colorLabel = "12"  // blue   — labels, host headers
	colorCyan  = "14"  // cyan   — sub-labels, usernames
	colorPort  = "11"  // light cyan — host port
	colorTree  = "212" // pink   — tree branch characters
)

// Brand returns the brand pill style (reverse video) used for the
// MINOP prompt label.
func Brand() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(colorBrand)).
		Bold(true)
}

// PromptArrow is the dim "›" separator between the brand pill and
// the input field in the interactive REPL.
func PromptArrow() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBrand)).
		Bold(true)
}

// Title is used for section headers ("MINOP CLI COMMANDS",
// version banner).
func Title() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBrand)).
		Bold(true)
}

// Task is used for task names and the verb of a help entry.
func Task() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTask))
}

// HelpKey is the verb/key column of a help entry, distinct from Task
// only by intent; both share the same colour.
func HelpKey() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTask))
}

// Dim is the faint colour used for separators, time stamps, and the
// description column of help entries.
func Dim() lipgloss.Style {
	return lipgloss.NewStyle().
		Faint(true)
}

// Label is used for inline labels next to values (e.g. "• Version:
//
//	<value>").
func Label() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorLabel))
}

// HostHeader is the bold role-name header in `minop host` (e.g.
// "• all").
func HostHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Faint(false).
		Foreground(lipgloss.Color(colorLabel))
}

// HostUser is the user segment in `minop host` tree output.
func HostUser() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorCyan))
}

// HostAddr is the address segment in `minop host` tree output.
func HostAddr() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTask))
}

// HostPort is the port segment in `minop host` tree output.
func HostPort() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorPort))
}

// HostSep is the dim "@" and ":" separators between user, address,
// and port in `minop host`.
func HostSep() lipgloss.Style {
	return lipgloss.NewStyle().
		Faint(true)
}

// TreeBranch is the colour of tree branch characters (├──, └──).
func TreeBranch() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTree))
}

// ResultLabel is the colour of per-key labels inside a module's
// result output (e.g. "stdout:", "stderr:", "ExitStatus:" in the
// key/value block printed under each host line in `minop run`).
func ResultLabel() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorCyan))
}

// HostLine is the colour of the per-host line printed under each
// module header in `minop run` output.
func HostLine() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorLabel))
}

// Timestamp is the faint timestamp printed next to each host result
// in `minop run` output.
func Timestamp() lipgloss.Style {
	return lipgloss.NewStyle().
		Faint(true)
}
