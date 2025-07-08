package grpc_server

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "proto/smtppb"
	"results/errs"
	"results/succ"
	"smtp/internal/send_mail"
)

type Server struct {
	pb.UnimplementedSmtpServiceServer
}

func (s *Server) SendOtp(_ context.Context, req *pb.SendOtpRequest) (*pb.BaseResultResponse, error) {
	email := req.GetEmail()
	code := req.GetCode()

	if code == 0 || email == "" {
		return &pb.BaseResultResponse{Result: errs.SmtpCodeOrEmailMissing.Error()}, errs.SmtpCodeOrEmailMissing
	}

	if err := send_mail.SendOtp(email, code); err != nil {
		return &pb.BaseResultResponse{Result: err.Error()}, err
	}

	return &pb.BaseResultResponse{Result: succ.Ok}, nil
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.BaseResultResponse, error) {
	return &pb.BaseResultResponse{Result: succ.Ok}, nil
}
