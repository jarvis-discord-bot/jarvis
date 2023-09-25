package common

type InternalServerError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	EventID string `json:"event_id,omitempty"`
}

func NewInternalServerError(message string, eventID string) error {
	return &InternalServerError{
		Code:    "100001",
		Message: message,
		EventID: eventID,
	}
}

func (ise *InternalServerError) Error() string {
	return "error_code: " + ise.Code + " message: " + ise.Message
}

type BadRequestError struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func NewBadRequestError(message string, details interface{}) error {
	return &BadRequestError{
		Code:    "100002",
		Message: message,
		Details: details,
	}
}

func NewBadRequestErrorWithoutDetails(message string) error {
	return &BadRequestError{
		Code:    "100003",
		Message: message,
	}
}

func (bde *BadRequestError) Error() string {
	return "error_code: " + bde.Code + " message: " + bde.Message
}

type ConflictError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewConflictError(message string) error {
	return &ConflictError{
		Code:    "100004",
		Message: message,
	}
}

func (ce *ConflictError) Error() string {
	return "error_code: " + ce.Code + " message: " + ce.Message
}

type NotFoundError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewNotFoundError(id string) error {
	return &NotFoundError{
		Code:    "100005",
		Message: "path: " + id + " not found",
	}
}

func (nfe *NotFoundError) Error() string {
	return "error_code: " + nfe.Code + " message: " + nfe.Message
}

type ExpiredLinkError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewExpiredLinkError(id string, expirationDate string) error {
	return &ExpiredLinkError{
		Code:    "100006",
		Message: "id: " + id + " has expired at " + expirationDate,
	}
}

func (ele *ExpiredLinkError) Error() string {
	return "error_code: " + ele.Code + " message: " + ele.Message
}
