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

package action

import (
	"errors"
	"fmt"
)

var ErrParameterMissing = errors.New("required parameter is missing")

type Action interface {
	Validate(actCtx map[string]string) error
	Execute() (map[string]string, error)
}

func GetActionParam(actCtx map[string]string, key string) (string, error) {
	value, ok := actCtx[key]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrParameterMissing, key)
	}
	return value, nil
}
