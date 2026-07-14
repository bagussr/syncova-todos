package domain

type BaseResponse struct {
	Success    bool   `json:"success" default:"true"`
	StatusCode int    `json:"status_code" default:"200"`
	Message    string `json:"message" default:""`
}

type BaseListResponse struct {
	Success    bool   `json:"success" default:"true"`
	StatusCode int    `json:"status_code" default:"200"`
	Message    string `json:"message" default:""`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Total      int    `json:"total"`
}

func NewBaseResponse(success bool, statusCode int, message string) BaseResponse {
	return BaseResponse{
		Success:    success,
		StatusCode: statusCode,
		Message:    message,
	}
}
