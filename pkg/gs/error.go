package gs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrCanceled = status.New(codes.Canceled, "canceled").Err()
