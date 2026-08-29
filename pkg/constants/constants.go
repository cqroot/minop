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

package constants

// RoleAll is the special role that matches all host groups.
const RoleAll = "all"

const DefaultTaskFile = "minop.yaml"

// DefaultMaxProcs is the default maximum number of operations executed
// concurrently per task. It is a sensible default for small to medium
// fleets (10-50 hosts); users with larger fleets or tighter latency
// budgets can override it via the --max-procs flag.
const DefaultMaxProcs = 10
