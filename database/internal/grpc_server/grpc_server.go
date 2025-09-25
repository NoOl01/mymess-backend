package grpc_server

import (
	"context"
	"database/internal/bcrypt"
	"database/internal/common"
	"database/internal/db_models"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	pb "proto/databasepb"
	"results/errs"
	"results/succ"
	"strconv"
	"strings"
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
		err = common.CheckRecord(s.Db, "username", login.Username, &user)
	case *pb.LoginUserRequest_Email:
		err = common.CheckRecord(s.Db, "email", login.Email, &user)
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

func (s *Server) UpdatePassword(_ context.Context, req *pb.UpdatePasswordRequest) (*pb.BaseResultResponse, error) {
	email := req.GetEmail()
	pass := req.GetPassword()

	var user db_models.User
	if err := common.CheckRecord(s.Db, "email", email, &user); err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	hash, err := bcrypt.Encrypt(pass)
	if err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	if err := s.Db.Model(&user).Update("password", hash).Error; err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	return &pb.BaseResultResponse{
		Result: succ.Ok,
	}, nil
}

func (s *Server) UpdateProfile(_ context.Context, req *pb.UpdateRequest) (*pb.BaseResultResponse, error) {
	value := req.GetValue()
	updateType := req.GetType()
	userId := req.GetUserId()

	var user db_models.User
	if err := s.Db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.BaseResultResponse{
				Result: errs.RecordNotFound.Error(),
			}, errs.RecordNotFound
		}
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	switch updateType {
	case "nickname":
		user.Nickname = value
	case "email":
		user.Email = value
	case "avatar":
		user.AvatarPath = value
	case "banner":
		user.BannerPath = value
	case "bio":
		user.Bio = value
	default:
		return &pb.BaseResultResponse{
			Result: errs.UnknownType.Error(),
		}, errs.UnknownType
	}

	if err := s.Db.Save(&user).Error; err != nil {
		return &pb.BaseResultResponse{
			Result: err.Error(),
		}, err
	}

	return &pb.BaseResultResponse{
		Result: succ.Ok,
	}, nil
}

func (s *Server) FindProfile(_ context.Context, req *pb.FindProfileRequest) (*pb.FindProfileResponse, error) {
	value := req.GetName()
	var users []db_models.User
	var nameType string

	if strings.HasPrefix(value, "@") {
		nameType = "username"
		value = strings.TrimPrefix(value, "@")
	} else {
		nameType = "nickname"
	}

	if err := findProfile(fmt.Sprintf("%s ILIKE ?", nameType), value, &users, s.Db); err != nil {
		return &pb.FindProfileResponse{
			Error: err.Error(),
		}, err
	}

	var bodies []*pb.FindProfileBody
	for _, user := range users {
		bodies = append(bodies, &pb.FindProfileBody{
			Id:           user.Id,
			Nickname:     user.Nickname,
			Username:     user.Username,
			Avatar:       user.AvatarPath,
			Banner:       user.BannerPath,
			RegisteredAt: user.RegisterDate.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &pb.FindProfileResponse{
		Body:  bodies,
		Error: "",
	}, nil
}

func (s *Server) GetProfileInfo(_ context.Context, req *pb.GetProfileInfoRequest) (*pb.GetProfileInfoResponse, error) {
	id := req.GetId()
	var user db_models.User
	if err := s.Db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.GetProfileInfoResponse{
				Error: errs.RecordNotFound.Error(),
			}, errs.RecordNotFound
		}
		return &pb.GetProfileInfoResponse{
			Error: err.Error(),
		}, err
	}

	body := pb.FindProfileBody{
		Id:           user.Id,
		Nickname:     user.Nickname,
		Username:     user.Username,
		Avatar:       user.AvatarPath,
		Banner:       user.BannerPath,
		RegisteredAt: user.RegisterDate.Format("2006-01-02T15:04:05Z"),
	}

	return &pb.GetProfileInfoResponse{
		Body: &body,
	}, nil
}

func (s *Server) MyProfile(_ context.Context, req *pb.GetProfileInfoRequest) (*pb.MyProfileResponse, error) {
	id := req.GetId()

	var user db_models.User
	if err := s.Db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.MyProfileResponse{
				Error: errs.RecordNotFound.Error(),
			}, errs.RecordNotFound
		}
		return &pb.MyProfileResponse{
			Error: err.Error(),
		}, err
	}

	body := pb.FindProfileBody{
		Id:           user.Id,
		Nickname:     user.Nickname,
		Username:     user.Username,
		Avatar:       user.AvatarPath,
		Banner:       user.BannerPath,
		RegisteredAt: user.RegisterDate.Format("2006-01-02T15:04:05Z"),
	}

	return &pb.MyProfileResponse{
		Body: &body,
	}, nil
}

func findProfile(where, username string, users *[]db_models.User, db *gorm.DB) error {
	if err := db.Where(where, "%"+username+"%").Limit(10).Find(users).Error; err != nil {
		return err
	}
	return nil
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.BaseResultResponse, error) {
	return &pb.BaseResultResponse{Result: succ.Ok}, nil
}
