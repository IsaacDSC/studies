package main

import (
	"regexp"
	"strconv"
	"strings"
)

var intArgPattern = regexp.MustCompile(`^(min|max)=(\-?\d+)$`)
var floatArgPattern = regexp.MustCompile(`^(min|max)=(-?\d+(?:\.\d+)?)$`)
var funcCallPattern = regexp.MustCompile(`^([A-Za-z]+)\((.*)\)$`)

func parseFunctionRule(rule string) (string, []string, bool) {
	m := funcCallPattern.FindStringSubmatch(strings.TrimSpace(rule))
	if len(m) != 3 {
		return "", nil, false
	}
	name := strings.ToLower(strings.TrimSpace(m[1]))
	argText := strings.TrimSpace(m[2])
	if argText == "" {
		return name, nil, true
	}

	if strings.Contains(argText, "||") {
		items := strings.Split(argText, "||")
		return name, normalizeArgs(items), true
	}
	if strings.Contains(argText, "|") {
		items := strings.Split(argText, "|")
		return name, normalizeArgs(items), true
	}
	if strings.Contains(argText, ",") {
		items := strings.Split(argText, ",")
		return name, normalizeArgs(items), true
	}

	return name, normalizeArgs([]string{argText}), true
}

func normalizeArgs(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		s := strings.TrimSpace(item)
		s = strings.Trim(s, "\"'")
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func quoteArgs(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, strconv.Quote(arg))
	}
	return strings.Join(quoted, ", ")
}

func splitPipeList(value string) []string {
	items := strings.Split(value, "|")
	return normalizeArgs(items)
}
