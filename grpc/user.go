package grpc

import (
	"context"
	userv1 "github.com/MuxiKeStack/be-api/gen/proto/user/v1"
	"github.com/MuxiKeStack/be-user/domain"
	"github.com/MuxiKeStack/be-user/service"
	"google.golang.org/grpc"
)

type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
	svc service.UserService
}

func NewUserServiceServer(svc service.UserService) *UserServiceServer {
	return &UserServiceServer{svc: svc}
}

func (s *UserServiceServer) Register(server grpc.ServiceRegistrar) {
	userv1.RegisterUserServiceServer(server, s)
}

func (s *UserServiceServer) LoginByCCNU(ctx context.Context, request *userv1.LoginByCCNURequest) (*userv1.LoginByCCNUResponse, error) {
	u, err := s.svc.LoginByCCNU(ctx, request.StudentId, request.Password)
	switch err {
	case nil:
		return &userv1.LoginByCCNUResponse{User: convertToV(u)}, nil
	case service.ErrInvalidStudentIdOrPassword:
		return nil, userv1.ErrorInvalidSidOrPwd("学号或密码错误")
	default:
		return nil, err
	}
}

func (s *UserServiceServer) UpdateNonSensitiveInfo(ctx context.Context, request *userv1.UpdateNonSensitiveInfoRequest) (*userv1.UpdateNonSensitiveInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *UserServiceServer) Profile(ctx context.Context, request *userv1.ProfileRequest) (*userv1.ProfileResponse, error) {
	//TODO implement me
	panic("implement me")
}

func convertToV(user domain.User) *userv1.User {
	return &userv1.User{
		Id:        user.Id,
		StudentId: user.StudentId,
		Password:  user.Password,
		Avatar:    user.Avatar,
		Nickname:  user.Nickname,
		Utime:     user.Utime.UnixMilli(),
		Ctime:     user.Ctime.UnixMilli(),
		New:       user.New,
	}
}
