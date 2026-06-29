package domain

type BaseResponse[T any] struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

type BaseListResponse[T any] struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Total      int    `json:"total"`
	Data       []T    `json:"data"`
}

func NewSuccessResponse[T any](data T, message string) BaseResponse[T] {
	return BaseResponse[T]{
		Success:    true,
		StatusCode: 200,
		Message:    message,
		Data:       data,
	}
}

func NewCreatedResponse[T any](data T, message string) BaseResponse[T] {
	return BaseResponse[T]{
		Success:    true,
		StatusCode: 201,
		Message:    message,
		Data:       data,
	}
}

func NewErrorResponse(statusCode int, message string) BaseResponse[any] {
	return BaseResponse[any]{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
	}
}

func NewListResponse[T any](data []T, page, perPage, total int, message string) BaseListResponse[T] {
	return BaseListResponse[T]{
		Success:    true,
		StatusCode: 200,
		Message:    message,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		Data:       data,
	}
}
