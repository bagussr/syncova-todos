package domain

type BaseResponse struct {
	Success    bool   `json:"success" default:"true" binding:"required"`
	StatusCode int    `json:"status_code" default:"200" binding:"required"`
	Message    string `json:"message" binding:"required"`
}

type BaseListResponse struct {
	Success    bool        `json:"success" default:"true" binding:"required"`
	StatusCode int         `json:"status_code" default:"200" binding:"required"`
	Message    string      `json:"message" binding:"required"`
	Page       int         `json:"page" default:"1" binding:"required"`
	PerPage    int         `json:"per_page" default:"10" binding:"required"`
	Total      int64       `json:"total" default:"0" binding:"required"`
	Data       interface{} `json:"data" binding:"required"`
}

func NewBaseResponse(success bool, statusCode int, message string) BaseResponse {
	return BaseResponse{
		Success:    success,
		StatusCode: statusCode,
		Message:    message,
	}
}
