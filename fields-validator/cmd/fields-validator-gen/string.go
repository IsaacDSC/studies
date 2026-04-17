package main

import (
	"fmt"
	"strconv"
	"strings"
)

func mapStringRule(rule string) string {
	if fn, args, ok := parseFunctionRule(rule); ok {
		switch fn {
		case "equals":
			if len(args) != 1 {
				fail("equals() requires exactly one argument")
			}
			return fmt.Sprintf("validate.Equals(%s)", strconv.Quote(args[0]))
		case "contains":
			if len(args) != 1 {
				fail("contains() requires exactly one argument")
			}
			return fmt.Sprintf("validate.Contains(%s)", strconv.Quote(args[0]))
		case "oneof", "possibilities":
			if len(args) == 0 {
				fail("%s() requires at least one argument", fn)
			}
			return "validate.OneOf(" + quoteArgs(args) + ")"
		}
	}

	switch rule {
	case "nonEmpty":
		return "validate.Required()"
	case "email":
		return "validate.Email()"
	}

	if strings.HasPrefix(rule, "oneOf=") {
		items := splitPipeList(strings.TrimPrefix(rule, "oneOf="))
		return "validate.OneOf(" + quoteArgs(items) + ")"
	}

	matches := intArgPattern.FindStringSubmatch(rule)
	if len(matches) == 3 {
		n, err := strconv.Atoi(matches[2])
		if err != nil {
			fail("invalid numeric rule %q", rule)
		}
		if matches[1] == "min" {
			return fmt.Sprintf("validate.MinLen(%d)", n)
		}
		return fmt.Sprintf("validate.MaxLen(%d)", n)
	}

	fail("unsupported string rule %q", rule)
	return ""
}
