package validate

import (
	"strconv"
	"strings"
)

type IntRule func(field string, value int) []FieldError
type FloatRule func(field string, value float64) []FieldError

func ApplyIntRules(field string, value int, rules ...IntRule) []FieldError {
	var errs []FieldError
	for _, rule := range rules {
		errs = append(errs, rule(field, value)...)
	}
	return errs
}

func ApplyFloatRules(field string, value float64, rules ...FloatRule) []FieldError {
	var errs []FieldError
	for _, rule := range rules {
		errs = append(errs, rule(field, value)...)
	}
	return errs
}

func Min(min int) IntRule {
	return func(field string, value int) []FieldError {
		if value < min {
			return []FieldError{{
				Field: field, Code: "min", Message: "value is below minimum", Value: value,
			}}
		}
		return nil
	}
}

func Max(max int) IntRule {
	return func(field string, value int) []FieldError {
		if value > max {
			return []FieldError{{
				Field: field, Code: "max", Message: "value exceeds maximum", Value: value,
			}}
		}
		return nil
	}
}

func EqualsInt(expected int) IntRule {
	return func(field string, value int) []FieldError {
		if value != expected {
			return []FieldError{{
				Field: field, Code: "equals", Message: "value does not match expected", Value: value,
			}}
		}
		return nil
	}
}

func ContainsInt(expected int) IntRule {
	expectedToken := strconv.Itoa(expected)
	return func(field string, value int) []FieldError {
		if !strings.Contains(strconv.Itoa(value), expectedToken) {
			return []FieldError{{
				Field: field, Code: "contains", Message: "value does not contain expected token", Value: value,
			}}
		}
		return nil
	}
}

func OneOfInt(options ...int) IntRule {
	set := make(map[int]struct{}, len(options))
	for _, option := range options {
		set[option] = struct{}{}
	}
	return func(field string, value int) []FieldError {
		if _, ok := set[value]; !ok {
			return []FieldError{{
				Field: field, Code: "one_of", Message: "value is not in allowed options", Value: value,
			}}
		}
		return nil
	}
}

func MinFloat(min float64) FloatRule {
	return func(field string, value float64) []FieldError {
		if value < min {
			return []FieldError{{
				Field: field, Code: "min", Message: "value is below minimum", Value: value,
			}}
		}
		return nil
	}
}

func MaxFloat(max float64) FloatRule {
	return func(field string, value float64) []FieldError {
		if value > max {
			return []FieldError{{
				Field: field, Code: "max", Message: "value exceeds maximum", Value: value,
			}}
		}
		return nil
	}
}

func EqualsFloat(expected float64) FloatRule {
	return func(field string, value float64) []FieldError {
		if value != expected {
			return []FieldError{{
				Field: field, Code: "equals", Message: "value does not match expected", Value: value,
			}}
		}
		return nil
	}
}

func ContainsFloat(expected float64) FloatRule {
	expectedToken := strconv.FormatFloat(expected, 'f', -1, 64)
	return func(field string, value float64) []FieldError {
		valueToken := strconv.FormatFloat(value, 'f', -1, 64)
		if !strings.Contains(valueToken, expectedToken) {
			return []FieldError{{
				Field: field, Code: "contains", Message: "value does not contain expected token", Value: value,
			}}
		}
		return nil
	}
}

func OneOfFloat(options ...float64) FloatRule {
	set := make(map[float64]struct{}, len(options))
	for _, option := range options {
		set[option] = struct{}{}
	}
	return func(field string, value float64) []FieldError {
		if _, ok := set[value]; !ok {
			return []FieldError{{
				Field: field, Code: "one_of", Message: "value is not in allowed options", Value: value,
			}}
		}
		return nil
	}
}
