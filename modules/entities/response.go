package entities

import (
	"github.com/Montheankul-K/E-Commerce-Application-Backend/packages/logger"
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, traceId, msg string) IResponse
	Res() error
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    *fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	TraceId string `json:"trace_id"` // Id ของ error นั้นๆ
	Msg     string `json:"message"`
}

func NewResponse(c *fiber.Ctx) IResponse {
	return &Response{
		Context: c,
	}
}

func (r *Response) Success(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	logger.InitLogger(r.Context, &r.Data).Print()
	return r
}

func (r *Response) Error(code int, traceId, msg string) IResponse {
	r.StatusCode = code
	r.ErrorRes = &ErrorResponse{
		TraceId: traceId,
		Msg:     msg,
	}
	r.IsError = true // แต่เดิมค่าเป็น false
	logger.InitLogger(r.Context, &r.Data).Print()
	return r
}

func (r *Response) Res() error {
	return r.Context.Status(r.StatusCode).JSON(func() any {
		if r.IsError {
			return &r.ErrorRes // return เป็น reference เพราะ pass มาเป็น pointer
		}
		return &r.Data
	}())
}

type PaginateRes struct {
	Data      any `json:"data"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalItem int `json:"total_item"`
}
