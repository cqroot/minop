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

package module

// commonModule holds the Name and Group fields that every module
// shares, plus the trivial getter/setter pair that the Module
// interface requires. Concrete module types (Shell, Local, Copy)
// embed it to inherit the common surface area.
type commonModule struct {
	name  string
	group string
}

// Name returns the module's display name.
func (m commonModule) Name() string { return m.name }

// SetName replaces the module's display name.
func (m *commonModule) SetName(name string) { m.name = name }

// Group returns the host group this module targets.
func (m commonModule) Group() string { return m.group }

// SetGroup replaces the target host group.
func (m *commonModule) SetGroup(group string) { m.group = group }
