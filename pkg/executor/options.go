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

import (
	"os"

	"github.com/cqroot/minop/pkg/logs"
)

// Option configures an Executor.
type Option func(e *Executor)

// WithVerboseLevel sets the verbosity level for executor output.
func WithVerboseLevel(verboseLevel int) Option {
	return func(e *Executor) {
		e.optVerboseLevel = verboseLevel
	}
}

// WithPrinter sets a custom Printer for rendering task, host, and
// key/value output. When not supplied, New builds a Printer that
// writes to os.Stdout using the verbosity level set by
// WithVerboseLevel on the same Executor (or zero if none was set).
func WithPrinter(p *Printer) Option {
	return func(e *Executor) {
		e.optPrinter = p
	}
}

// defaultPrinter builds the Printer used when WithPrinter is not
// supplied. It writes to os.Stdout and uses the verbose level already
// stored on the Executor.
func defaultPrinter(verboseLevel int) *Printer {
	return NewPrinter(os.Stdout, verboseLevel)
}

// WithMaxProcs sets the maximum number of concurrent modules. A
// non-positive value is rejected: the executor keeps its current
// optMaxProcs (the default set by New) and logs a warning so the
// caller can see that their input had no effect.
func WithMaxProcs(maxProcs int) Option {
	return func(e *Executor) {
		if maxProcs > 0 {
			e.optMaxProcs = maxProcs
			return
		}
		logs.Logger().Warn().
			Int("requested", maxProcs).
			Int("kept", e.optMaxProcs).
			Msg("WithMaxProcs: non-positive value ignored, keeping current maxProcs")
	}
}
