package grpc_server

import (
	"auth/internal/grpc_client"
	"auth/internal/jwt"
	"auth/internal/otp"
	"auth/internal/reset"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "proto/authpb"
	"proto/databasepb"
	"proto/smtppb"
	"results/errs"
	"results/succ"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	Client     databasepb.DatabaseServiceClient
	SmtpClient smtppb.SmtpServiceClient
}

func (s *Server) Register(_ context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	nickname := req.GetNickname()
	email := req.GetEmail()
	password := req.GetPassword()

	if nickname == "" || email == "" || password == "" {
		return &pb.AuthResponse{
			Error: "Username, Email and Password is empty",
		}, nil
	}

	resp, err := grpc_client.CreateNewUser(s.Client, nickname, email, password)
	if err != nil {
		return nil, err
	}

	if resp.Result != succ.RecordCreated {
		return nil, fmt.Errorf(resp.Result)
	}

	token, err := jwt.GenerateToken(resp.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{AccessToken: token}, nil
}

func (s *Server) Login(_ context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	var dbUsername, dbEmail string

	switch login := req.LoginMethod.(type) {
	case *pb.LoginRequest_Email:
		dbEmail = login.Email
	case *pb.LoginRequest_Username:
		dbUsername = login.Username
	default:
		return &pb.AuthResponse{Error: errs.InvalidRequestBody.Error()}, errs.InvalidRequestBody
	}

	resp, err := grpc_client.CheckUser(s.Client, dbUsername, dbEmail, req.GetPassword())
	if err != nil {
		return &pb.AuthResponse{Error: err.Error()}, err
	}

	if resp.Result != succ.Ok {
		return &pb.AuthResponse{Error: resp.Result}, fmt.Errorf(resp.Result)
	}

	token, err := jwt.GenerateToken(resp.UserId)
	if err != nil {
		return &pb.AuthResponse{Error: err.Error()}, err
	}

	return &pb.AuthResponse{AccessToken: token}, nil
}

func (s *Server) Refresh(_ context.Context, req *pb.RefreshRequest) (*pb.AuthResponse, error) {
	token := req.GetAccessToken()

	if token == "" {
		return &pb.AuthResponse{
			Error: errs.MissingToken.Error(),
		}, errs.MissingToken
	}

	userId, err := jwt.ValidateJwt(token)
	if err != nil {
		return &pb.AuthResponse{
			Error: err.Error(),
		}, err
	}

	accessToken, err := jwt.GenerateToken(userId)
	if err != nil {
		return &pb.AuthResponse{
			Error: err.Error(),
		}, err
	}

	return &pb.AuthResponse{AccessToken: accessToken}, nil
}

func (s *Server) SendOtp(_ context.Context, req *pb.SendOtpRequest) (*pb.BaseResultResponse, error) {
	email := req.GetEmail()
	if email == "" {
		return &pb.BaseResultResponse{
			Result: errs.SmtpEmailMissing.Error(),
		}, errs.SmtpEmailMissing
	}

	code := otp.Generate()
	otp.StoreOTP(email, code)

	resp, err := grpc_client.SendOtp(s.SmtpClient, email, code)
	if err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	return &pb.BaseResultResponse{
		Result: resp.Result,
	}, nil
}

func (s *Server) ResetPassword(_ context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	email := req.GetEmail()
	code := req.GetCode()

	if email == "" || code == 0 {
		return &pb.ResetPasswordResponse{
			Result: errs.SmtpCodeOrEmailMissing.Error(),
		}, errs.SmtpCodeOrEmailMissing
	}

	ok := otp.VerifyOTP(email, code)
	if !ok {
		return &pb.ResetPasswordResponse{
			Result: errs.SmtpWrongOtpCode.Error(),
		}, errs.SmtpWrongOtpCode
	}

	resetToken := reset.StoreToken(email)

	return &pb.ResetPasswordResponse{
		Result:     succ.Ok,
		ResetToken: resetToken,
	}, nil
}

func (s *Server) UpdatePassword(_ context.Context, req *pb.UpdatePasswordRequest) (*pb.BaseResultResponse, error) {
	email := req.GetEmail()
	pass := req.GetPassword()
	resetToken := req.GetResetToken()

	if email == "" || pass == "" || resetToken == "" {
		return &pb.BaseResultResponse{
			Result: errs.InvalidRequestBody.Error(),
		}, errs.InvalidRequestBody
	}

	ok := reset.VerifyToken(email, resetToken)
	if !ok {
		return &pb.BaseResultResponse{
			Result: errs.InvalidToken.Error(),
		}, errs.InvalidToken
	}

	resp, err := grpc_client.UpdatePassword(s.Client, email, pass)
	if err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	return &pb.BaseResultResponse{
		Result: resp.Result,
	}, nil
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.BaseResultResponse, error) {
	return &pb.BaseResultResponse{Result: succ.Ok}, nil
}
