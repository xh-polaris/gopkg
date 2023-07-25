package errors

import (
	"github.com/xh-polaris/gopkg/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BizError struct {
	Code uint32
	Msg  string
}

func (e *BizError) ToGRPCError() error {
	return status.Error(codes.Code(e.Code), e.Msg)
}

func (e *BizError) Error() string {
	return util.JSONF(e)
}
