package main

import (
	"fmt"
	"strconv"
	"strings"
)

func mapIntRule(rule string) string {
	if fn, args, ok := parseFunctionRule(rule); ok {
		switch fn {
		case "equals":
			if len(args) != 1 {
				fail("equals() requires exactly one argument")
			}
			n, err := strconv.Atoi(args[0])
			if err != nil {
				fail("equals() int argument must be integer: %q", args[0])
			}
			return fmt.Sprintf("validate.EqualsInt(%d)", n)
		case "contains":
			if len(args) != 1 {
				fail("contains() requires exactly one argument")
			}
			n, err := strconv.Atoi(args[0])
			if err != nil {
				fail("contains() int argument must be integer: %q", args[0])
			}
			return fmt.Sprintf("validate.ContainsInt(%d)", n)
		case "oneof", "possibilities":
			if len(args) == 0 {
				fail("%s() requires at least one argument", fn)
			}
			intArgs := make([]string, 0, len(args))
			for _, arg := range args {
				n, err := strconv.Atoi(arg)
				if err != nil {
					fail("%s() int argument must be integer: %q", fn, arg)
				}
				intArgs = append(intArgs, strconv.Itoa(n))
			}
			return "validate.OneOfInt(" + strings.Join(intArgs, ", ") + ")"
		}
	}

	matches := intArgPattern.FindStringSubmatch(rule)
	if len(matches) != 3 {
		fail("unsupported int rule %q", rule)
	}
	n, err := strconv.Atoi(matches[2])
	if err != nil {
		fail("invalid numeric rule %q", rule)
	}
	if matches[1] == "min" {
		return fmt.Sprintf("validate.Min(%d)", n)
	}
	return fmt.Sprintf("validate.Max(%d)", n)
}

func mapIntPtrRule(rule string) string {
	r := strings.TrimSpace(rule)
	if strings.EqualFold(r, "notnil") || strings.EqualFold(r, "nonnil") {
		return "validate.PtrNotNilInt()"
	}
	return "validate.PtrInt(" + mapIntRule(rule) + ")"
}
