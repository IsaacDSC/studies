package validate

import (
	"strings"
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
