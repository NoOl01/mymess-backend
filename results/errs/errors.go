package errs

import "errors"

var (
	ServerError        = errors.New("grpc_server error")
	InvalidRequestBody = errors.New("invalid request body")
	FailedListen       = errors.New("failed to listen")
	FailedLoadEnvFile  = errors.New("failed to load .env file")
	WrongPassword      = errors.New("wrong password")
)

// gRPC Client
var (
	GrpcClientConnectFailed = errors.New("failed connect to the gRPC grpc_server")
	GrpcClientCloseFailed   = errors.New("the connection to the gRPC grpc_server could not be closed")
)

// Database
var (
	FailedDatabaseConnect = errors.New("failed to connect to the database")
	FailedCreateRecord    = errors.New("failed to create record")
	FailedReadRecord      = errors.New("failed to read the record")
	RecordNotFound        = errors.New("record not found")
	RecordAlreadyExists   = errors.New("record already exists")
)
