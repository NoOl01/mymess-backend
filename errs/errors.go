package errs

import "errors"

var (
	GrpcClientConnectFailed = errors.New("failed connect to the gRPC grpc_server")
	GrpcClientCloseFailed   = errors.New("the connection to the gRPC grpc_server could not be closed")
	ServerError             = errors.New("grpc_server error")
	AuthLoginInvalidBody    = errors.New("invalid request body")
	FailedListen            = errors.New("failed to listen")
	FailedServe             = errors.New("failed to serve")
)
