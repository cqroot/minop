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

package utils

import (
	"slices"
	"strconv"
	"strings"
)

func StrIsInteger(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func StrToInteger(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func StrIsFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func StrToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

var (
	trueStrings  = []string{"true", "yes", "y"}
	falseStrings = []string{"false", "no", "n"}
)

func StrIsBoolean(s string) bool {
	if slices.Contains(trueStrings, strings.ToLower(s)) || slices.Contains(falseStrings, strings.ToLower(s)) {
		return true
	}
	return false
}

func StrToBoolean(s string) bool {
	return slices.Contains(trueStrings, strings.ToLower(s))
}
