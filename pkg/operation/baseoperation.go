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

package operation

// baseOperation defines the interface for common operation properties.
type baseOperation interface {
	Name() string
	SetName(string)
	Role() string
	SetRole(string)
}

// baseOperationImpl provides a base implementation for operations.
type baseOperationImpl struct {
	name string
	role string
}

// Name returns the operation's name.
func (op baseOperationImpl) Name() string {
	return op.name
}

// SetName sets the operation's name.
func (op *baseOperationImpl) SetName(name string) {
	op.name = name
}

// Role returns the operation's target role.
func (op baseOperationImpl) Role() string {
	return op.role
}

// SetRole sets the operation's target role.
func (op *baseOperationImpl) SetRole(role string) {
	op.role = role
}
