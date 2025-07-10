package models

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewSuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
}

func NewErrorResponse(message string) APIResponse {
	return APIResponse{
		Success: false,
		Data:    nil,
		Message: message,
	}
}
