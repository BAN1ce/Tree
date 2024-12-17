package controller

import "github.com/BAN1ce/Tree/inner/api"

type BaseController struct {
	Api api.API
}

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

func NewResponse(success bool, code int, msg string, data interface{}) *Response {
	return &Response{
		Success: success,
		Code:    code,
		Data:    data,
		Message: msg,
	}
}
