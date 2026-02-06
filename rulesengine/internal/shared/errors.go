package shared

import "fmt"

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common errors
var (
	ErrInvalidUUID = &AppError{
		Code:    "invalid_user_id",
		Message: "user_id must be a valid UUID",
		Status:  400,
	}

	ErrInvalidOperation = &AppError{
		Code:    "invalid_operation",
		Message: "operation must be 'credit' or 'anticipation'",
		Status:  400,
	}

	ErrInvalidAmount = &AppError{
		Code:    "invalid_amount",
		Message: "requested_amount must be a positive integer (cents)",
		Status:  400,
	}

	ErrUserNotFound = &AppError{
		Code:    "user_not_found",
		Message: "user not found in the system",
		Status:  422,
	}

	ErrUserAlreadyBlocked = &AppError{
		Code:    "user_already_blocked",
		Message: "user is already blocked",
		Status:  409,
	}

	ErrUserNotBlocked = &AppError{
		Code:    "user_not_blocked",
		Message: "user is not blocked",
		Status:  404,
	}

	ErrCampaignNotFound = &AppError{
		Code:    "campaign_not_found",
		Message: "campaign not found",
		Status:  404,
	}

	ErrNoCampaigns = &AppError{
		Code:    "no_campaigns",
		Message: "no campaigns configured for this operation",
		Status:  422,
	}

	ErrInvalidGraph = &AppError{
		Code:    "invalid_graph",
		Message: "invalid rules graph",
		Status:  422,
	}

	ErrInvalidJSON = &AppError{
		Code:    "invalid_json",
		Message: "invalid JSON body",
		Status:  400,
	}

	ErrMissingOperations = &AppError{
		Code:    "missing_operations",
		Message: "operations query parameter is required",
		Status:  400,
	}
)

// NewAppError creates a new application error
func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}


