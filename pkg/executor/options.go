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

package executor

// Option configures an Executor.
type Option func(e *Executor)

// WithVerboseLevel sets the verbosity level for executor output.
func WithVerboseLevel(verboseLevel int) Option {
	return func(e *Executor) {
		e.optVerboseLevel = verboseLevel
	}
}

// WithMaxProcs sets the maximum number of concurrent operations.
// A value of 0 or negative is ignored and the default (1) is used.
func WithMaxProcs(maxProcs int) Option {
	return func(e *Executor) {
		if maxProcs > 0 {
			e.optMaxProcs = maxProcs
		}
	}
}
