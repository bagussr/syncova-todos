package dto

import "syncova-todo/domain"

type TodoResponse struct {
	domain.BaseResponse `json:",inline"`
	Data                string `json:"data"`
}

func NewTodoSuccessResponse(data string, message string) TodoResponse {
	return TodoResponse{
		BaseResponse: domain.BaseResponse{
			Success:    true,
			StatusCode: 200,
			Message:    message,
		},
		Data: data,
	}
}
