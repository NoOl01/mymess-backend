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

// JWT
var (
	InvalidToken           = errors.New("invalid token")
	InvalidTokenClaimsType = errors.New("invalid claims type")
	InvalidOrMissingClaim  = errors.New("invalid or missing claim")
	UnexpectedSignMethod   = errors.New("unexpected signing method")
	MissingToken           = errors.New("missing token")
)

// WebSocket
var (
	WSClientCloseFailed    = errors.New("failed to close client connection")
	WsUpgradeFailed        = errors.New("failed to upgrade websocket")
	WsDecodeJsonFailed     = errors.New("failed to decode json")
	WsListenAndServeFailed = errors.New("failed to listen and serve")
)

// Smtp
var (
	SmtpCodeOrEmailMissing = errors.New("missing code or email")
	SmtpEmailMissing       = errors.New("email missing")
	SmtpWrongOtpCode       = errors.New("wrong otp code")
)
