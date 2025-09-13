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

package maputils

import (
	"errors"
	"fmt"

	"github.com/cqroot/minop/pkg/utils/strutils"
)

var (
	ErrKeyNotFound = errors.New("key not found in map")
	ErrValueType   = errors.New("error value type in map")
)

func GetString(strMap map[string]string, key string) (string, error) {
	value, ok := strMap[key]
	if !ok {
		return "", fmt.Errorf("%w: %s (%+v)", ErrKeyNotFound, key, strMap)
	}
	return value, nil
}

func GetStringOrDefault(strMap map[string]string, key string, def string) string {
	value, err := GetString(strMap, key)
	if err != nil {
		return def
	}
	return value
}

func GetBool(strMap map[string]string, key string) (bool, error) {
	value, ok := strMap[key]
	if !ok {
		return false, fmt.Errorf("%w: %s (%+v)", ErrKeyNotFound, key, strMap)
	}
	if !strutils.IsBool(value) {
		return false, fmt.Errorf("%w: bool", ErrValueType)
	}
	return strutils.ToBool(value), nil
}

func GetBoolOrDefault(strMap map[string]string, key string, def bool) bool {
	value, err := GetBool(strMap, key)
	if err != nil {
		return def
	}
	return value
}
