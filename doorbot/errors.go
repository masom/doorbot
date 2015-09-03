package doorbot

// ErrorResponse is the base class for API error responses
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// BadRequestErrorResponse HTTP 400
type BadRequestErrorResponse struct {
	ErrorResponse
}

// UnauthorizedErrorResponse HTTP 401
type UnauthorizedErrorResponse struct {
	ErrorResponse
}

// ForbiddenErrorResponse HTTP 403
type ForbiddenErrorResponse struct {
	ErrorResponse
}

// ConflictErrorResponse HTTP 409
type ConflictErrorResponse struct {
	ErrorResponse
}

// EntityNotFoundResponse HTTP 404
type EntityNotFoundResponse struct {
	ErrorResponse
}

// InternalServerErrorReponse HTTP 500
type InternalServerErrorReponse struct {
	ErrorResponse
}

// ServiceUnavailableErrorResponse HTTP 503
type ServiceUnavailableErrorResponse struct {
	ErrorResponse
}

// NewConflictErrorResponse creates a new ConflictErrorResponse
func NewConflictErrorResponse(errors []string) *ConflictErrorResponse {
	response := &ConflictErrorResponse{}
	response.Errors = errors
	return response
}

// NewBadRequestErrorResponse creates a new BadRequestErrorResponse
func NewBadRequestErrorResponse(errors []string) *BadRequestErrorResponse {
	response := &BadRequestErrorResponse{}
	response.Errors = errors
	return response
}

// NewEntityNotFoundResponse creates a new EntityNotFoundResponse
func NewEntityNotFoundResponse(errors []string) *EntityNotFoundResponse {
	response := &EntityNotFoundResponse{}
	response.Errors = errors
	return response
}

// NewInternalServerErrorResponse creates a new InternalServerErrorResponse
func NewInternalServerErrorResponse(errors []string) *InternalServerErrorReponse {
	response := &InternalServerErrorReponse{}
	response.Errors = errors
	return response
}

// NewServiceUnavailableErrorResponse creates a new ServiceUnavailableErrorResponse
func NewServiceUnavailableErrorResponse(errors []string) *ServiceUnavailableErrorResponse {
	response := &ServiceUnavailableErrorResponse{}
	response.Errors = errors

	return response
}

// NewUnauthorizedErrorResponse creates a new UnauthorizedErrorResponse
func NewUnauthorizedErrorResponse(errors []string) *UnauthorizedErrorResponse {
	response := &UnauthorizedErrorResponse{}
	response.Errors = errors

	return response
}

// NewForbiddenErrorResponse creates a new ForbiddenErrorResponse
func NewForbiddenErrorResponse(errors []string) *ForbiddenErrorResponse {
	response := &ForbiddenErrorResponse{}
	response.Errors = errors

	return response
}
