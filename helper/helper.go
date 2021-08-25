package helper

import (
	"github.com/online-store/consts"
	"github.com/online-store/entity"
)

func NewResponse(code uint, message consts.ResponseMessage, data interface{}) entity.Response {
	response := entity.Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return response
}
