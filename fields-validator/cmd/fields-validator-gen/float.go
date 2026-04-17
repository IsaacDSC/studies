package main

import (
	"fmt"
	"strconv"
	"strings"
)

func mapFloatRule(rule string) string {
	if fn, args, ok := parseFunctionRule(rule); ok {
		switch fn {
		case "equals":
			if len(args) != 1 {
				fail("equals() requires exactly one argument")
			}
			n, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				fail("equals() float argument must be numeric: %q", args[0])
			}
			return fmt.Sprintf("validate.EqualsFloat(%s)", strconv.FormatFloat(n, 'f', -1, 64))
		case "contains":
			if len(args) != 1 {
				fail("contains() requires exactly one argument")
			}
			n, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				fail("contains() float argument must be numeric: %q", args[0])
			}
			return fmt.Sprintf("validate.ContainsFloat(%s)", strconv.FormatFloat(n, 'f', -1, 64))
		case "oneof", "possibilities":
			if len(args) == 0 {
				fail("%s() requires at least one argument", fn)
			}
			floatArgs := make([]string, 0, len(args))
			for _, arg := range args {
				n, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					fail("%s() float argument must be numeric: %q", fn, arg)
				}
				floatArgs = append(floatArgs, strconv.FormatFloat(n, 'f', -1, 64))
			}
			return "validate.OneOfFloat(" + strings.Join(floatArgs, ", ") + ")"
		}
	}

	matches := floatArgPattern.FindStringSubmatch(rule)
	if len(matches) != 3 {
		fail("unsupported float rule %q", rule)
	}
	n, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		fail("invalid numeric rule %q", rule)
	}
	if matches[1] == "min" {
		return fmt.Sprintf("validate.MinFloat(%s)", strconv.FormatFloat(n, 'f', -1, 64))
	}
	return fmt.Sprintf("validate.MaxFloat(%s)", strconv.FormatFloat(n, 'f', -1, 64))
}
