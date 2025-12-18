package biz

import (
	"context"
	"net/http"
)

type Response struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(data any) *Response {
	return &Response{
		Code: 0,
		Msg:  "操作成功",
		Data: data,
	}
}

func Fail(code int64, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
	}
}

func SpecificFail(err *SpecificError) *Response {
	return &Response{
		Code: 1,
		Msg:  err.Msg,
	}
}

type SpecificError struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (e *SpecificError) Error() string {
	return e.Msg
}

func NewSpecificError(code int64, msg string) *SpecificError {
	return &SpecificError{
		Code: code,
		Msg:  msg,
	}
}

func OkHandler(_ context.Context, v any) any {
	return Success(v)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		switch e := err.(type) {
		case *SpecificError:
			return http.StatusOK, SpecificFail(e)
		default:
			return http.StatusOK, Fail(CodeServerError.Code, err.Error())
		}
	}
}
