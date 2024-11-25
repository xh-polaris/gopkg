package errors

import (
	"github.com/xh-polaris/gopkg/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BizError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

func (e *BizError) ToGRPCError() error {
	return status.Error(codes.Code(e.Code), e.Msg)
}

func (e *BizError) Error() string {
	return util.JSONF(e)
}

func NewBizError(code uint32, msg string) error {
	return (&BizError{
		Code: code,
		Msg:  msg,
	}).ToGRPCError()
}
