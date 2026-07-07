package models

import internal "github.com/AgileExecutives/serverbase/internal/models"

// Re-export helper functions and types from internal models for external use.
func SuccessResponse(message string, data interface{}) internal.APIResponse {
	return internal.SuccessResponse(message, data)
}

func SuccessMessageResponse(message string) internal.APIResponse {
	return internal.SuccessMessageResponse(message)
}

func SuccessListResponse(data interface{}, page, limit, total int) internal.APIResponse {
	return internal.SuccessListResponse(data, page, limit, total)
}

func ErrorResponseFunc(message string, err string) internal.APIResponse {
	return internal.ErrorResponseFunc(message, err)
}

func SimpleErrorResponse(err string) internal.ErrorResponse {
	return internal.SimpleErrorResponse(err)
}

type APIResponse = internal.APIResponse
type ErrorResponse = internal.ErrorResponse
type ListResponse = internal.ListResponse
type PaginationResponse = internal.PaginationResponse
