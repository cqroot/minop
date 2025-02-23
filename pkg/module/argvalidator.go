package module

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type ArgType string

const (
	ArgTypeString  ArgType = "String"
	ArgTypeInteger ArgType = "Integer"
	ArgTypeFloat   ArgType = "Float"
	ArgTypeBoolean ArgType = "Boolean"
)

type Arg struct {
	Name     string
	Desc     string
	Type     ArgType
	Optional bool
	Default  string
}

var (
	ErrArgNotFound     = errors.New("arg not found")
	ErrArgTypeMismatch = errors.New("arg type mismatch")
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

var trueStrings = []string{"true", "yes", "y"}
var falseStrings = []string{"false", "no", "n"}

func StrIsBoolean(s string) bool {
	if slices.Contains(trueStrings, strings.ToLower(s)) || slices.Contains(falseStrings, strings.ToLower(s)) {
		return true
	}
	return false
}

func StrToBoolean(s string) bool {
	if slices.Contains(trueStrings, strings.ToLower(s)) {
		return true
	}
	return false
}

func ValidateArgs(argMap map[string]string, argTypes []Arg) error {
	for _, arg := range argTypes {
		// Only mandatory args require validation
		if arg.Optional {
			continue
		}

		val, ok := argMap[arg.Name]
		if !ok {
			return fmt.Errorf("%w: %s", ErrArgNotFound, arg.Name)
		}

		switch arg.Type {
		case ArgTypeString:
			continue

		case ArgTypeInteger:
			if !StrIsInteger(val) {
				return fmt.Errorf("%w: %s %s -> %s", ErrArgTypeMismatch, arg.Name, arg.Type, val)
			}

		case ArgTypeFloat:
			if !StrIsFloat(val) {
				return fmt.Errorf("%w: %s %s -> %s", ErrArgTypeMismatch, arg.Name, arg.Type, val)
			}

		case ArgTypeBoolean:
			if !StrIsBoolean(val) {
				return fmt.Errorf("%w: %s %s -> %s", ErrArgTypeMismatch, arg.Name, arg.Type, val)
			}
		}
	}

	return nil
}
