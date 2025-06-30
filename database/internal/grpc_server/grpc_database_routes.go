package grpc_server

import (
	"context"
	"database/internal/bcrypt"
	"database/internal/db_models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	pb "proto/databasepb"
	"results/errs"
	"results/succ"
	"strconv"
	"time"
)

type Server struct {
	pb.UnimplementedDatabaseServiceServer
	Db *gorm.DB
}

func (s *Server) Register(_ context.Context, req *pb.CreateNewUserRequest) (*pb.AuthResponse, error) {
	nickname := req.GetNickname()
	email := req.GetEmail()
	password := req.GetPassword()

	if nickname == "" || email == "" || password == "" {
		return &pb.AuthResponse{
			UserId: "",
			Result: errs.InvalidRequestBody.Error(),
		}, nil
	}

	var user db_models.User

	if err := s.Db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			hash, err := bcrypt.Encrypt(password)
			if err != nil {
				return &pb.AuthResponse{
					UserId: "",
					Result: err.Error(),
				}, err
			}

			user = db_models.User{
				Nickname:     nickname,
				Email:        email,
				Password:     hash,
				Deleted:      false,
				RegisterDate: time.Now(),
			}

			if err := s.Db.Create(&user).Error; err != nil {
				return &pb.AuthResponse{
					UserId: "",
					Result: errs.FailedCreateRecord.Error(),
				}, err
			}
			return &pb.AuthResponse{
				UserId: strconv.FormatInt(user.Id, 10),
				Result: succ.RecordCreated,
			}, nil
		}
		return &pb.AuthResponse{
			UserId: "",
			Result: errs.FailedReadRecord.Error(),
		}, err
	}

	return &pb.AuthResponse{
		UserId: strconv.FormatInt(user.Id, 10),
		Result: errs.RecordAlreadyExists.Error(),
	}, nil
}

func (s *Server) Login(_ context.Context, req *pb.LoginUserRequest) (*pb.AuthResponse, error) {
	var user db_models.User
	var err error

	switch login := req.LoginMethod.(type) {
	case *pb.LoginUserRequest_Username:
		err = checkRecord(s.Db, "username", login.Username, &user)
	case *pb.LoginUserRequest_Email:
		err = checkRecord(s.Db, "email", login.Email, &user)
	default:
		return &pb.AuthResponse{
			UserId: "",
			Result: errs.InvalidRequestBody.Error(),
		}, errs.InvalidRequestBody
	}

	if err != nil {
		if errors.Is(err, errs.RecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return &pb.AuthResponse{
			UserId: "",
			Result: errs.RecordNotFound.Error(),
		}, err
	}

	if err := bcrypt.ValidatePassword(user, req.GetPassword()); err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		UserId: strconv.FormatInt(user.Id, 10),
		Result: succ.Ok,
	}, nil
}

func checkRecord[T any](db *gorm.DB, modelName string, value string, model *T) error {
	if err := db.Where(fmt.Sprintf("%s = ?", modelName), value).First(model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.RecordNotFound
		}
		return err
	}
	return nil
}
