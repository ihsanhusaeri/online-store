package entity

import "github.com/online-store/consts"

type Response struct {
	Code    uint                   `json:"code"`
	Message consts.ResponseMessage `json:"message"`
	Data    interface{}            `json:"data"`
}
