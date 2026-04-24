package validate

import (
	"strings"
	"time"
)

type StringRule func(field, value string) []FieldError

func ApplyStringRules(field, value string, rules ...StringRule) []FieldError {
	var errs []FieldError
	for _, rule := range rules {
		errs = append(errs, rule(field, value)...)
	}
	return errs
}

type StringPtrRule func(field string, value *string) []FieldError

func ApplyStringPtrRules(field string, value *string, rules ...StringPtrRule) []FieldError {
	var errs []FieldError
	for _, rule := range rules {
		errs = append(errs, rule(field, value)...)
	}
	return errs
}

func PtrNotNil() StringPtrRule {
	return func(field string, value *string) []FieldError {
		if value == nil {
			return []FieldError{{
				Field: field, Code: "not_nil", Message: "field must not be null", Value: nil,
			}}
		}
		return nil
	}
}

func PtrString(rule StringRule) StringPtrRule {
	return func(field string, value *string) []FieldError {
		if value == nil {
			return nil
		}
		return rule(field, *value)
	}
}

func Required() StringRule {
	return func(field, value string) []FieldError {
		if strings.TrimSpace(value) == "" {
			return []FieldError{{
				Field: field, Code: "required", Message: "field is required", Value: value,
			}}
		}
		return nil
	}
}

func MinLen(min int) StringRule {
	return func(field, value string) []FieldError {
		if len(value) < min {
			return []FieldError{{
				Field: field, Code: "min", Message: "length is below minimum", Value: value,
			}}
		}
		return nil
	}
}

func MaxLen(max int) StringRule {
	return func(field, value string) []FieldError {
		if len(value) > max {
			return []FieldError{{
				Field: field, Code: "max", Message: "length exceeds maximum", Value: value,
			}}
		}
		return nil
	}
}

func Email() StringRule {
	return func(field, value string) []FieldError {
		parts := strings.Split(value, "@")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" || !strings.Contains(parts[1], ".") {
			return []FieldError{{
				Field: field, Code: "invalid_format", Message: "invalid email format", Value: value,
			}}
		}
		return nil
	}
}

func Equals(expected string) StringRule {
	return func(field, value string) []FieldError {
		if value != expected {
			return []FieldError{{
				Field: field, Code: "equals", Message: "value does not match expected", Value: value,
			}}
		}
		return nil
	}
}

func Contains(expected string) StringRule {
	return func(field, value string) []FieldError {
		if !strings.Contains(value, expected) {
			return []FieldError{{
				Field: field, Code: "contains", Message: "value does not contain expected token", Value: value,
			}}
		}
		return nil
	}
}

func OneOf(options ...string) StringRule {
	set := make(map[string]struct{}, len(options))
	for _, option := range options {
		set[option] = struct{}{}
	}
	return func(field, value string) []FieldError {
		if _, ok := set[value]; !ok {
			return []FieldError{{
				Field: field, Code: "one_of", Message: "value is not in allowed options", Value: value,
			}}
		}
		return nil
	}
}

func Date(format string) StringRule {
	return func(field, value string) []FieldError {
		layout, ok := parseDateTimeLayout(format, "date")
		if !ok {
			return []FieldError{{
				Field: field, Code: "invalid_rule", Message: "unsupported date format rule", Value: format,
			}}
		}
		if _, err := time.Parse(layout, value); err != nil {
			return []FieldError{{
				Field: field, Code: "invalid_format", Message: "invalid date format", Value: value,
			}}
		}
		return nil
	}
}

func Time(format string) StringRule {
	return func(field, value string) []FieldError {
		layout, ok := parseDateTimeLayout(format, "time")
		if !ok {
			return []FieldError{{
				Field: field, Code: "invalid_rule", Message: "unsupported time format rule", Value: format,
			}}
		}
		if _, err := time.Parse(layout, value); err != nil {
			return []FieldError{{
				Field: field, Code: "invalid_format", Message: "invalid time format", Value: value,
			}}
		}
		return nil
	}
}

func DateTime(format string) StringRule {
	return func(field, value string) []FieldError {
		layout, ok := parseDateTimeLayout(format, "datetime")
		if !ok {
			return []FieldError{{
				Field: field, Code: "invalid_rule", Message: "unsupported datetime format rule", Value: format,
			}}
		}
		if _, err := time.Parse(layout, value); err != nil {
			return []FieldError{{
				Field: field, Code: "invalid_format", Message: "invalid datetime format", Value: value,
			}}
		}
		return nil
	}
}

func parseDateTimeLayout(format, kind string) (string, bool) {
	f := strings.TrimSpace(strings.ToUpper(format))
	switch kind {
	case "date":
		if f == "YYYY-MM-DD" {
			return "2006-01-02", true
		}
	case "time":
		if f == "HH:MM:SS" {
			return "15:04:05", true
		}
	case "datetime":
		if f == "YYYY-MM-DD HH:MM:SS" {
			return "2006-01-02 15:04:05", true
		}
	}
	return "", false
}
