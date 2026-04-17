package validate

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

type Errors []FieldError

func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}
	return "validation failed"
}

func (e Errors) HasAny() bool {
	return len(e) > 0
}
