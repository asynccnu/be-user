package grpc

import (
	"context"
	userv1 "github.com/asynccnu/be-api/gen/proto/user/v1"
	"github.com/asynccnu/be-user/domain"
	"github.com/asynccnu/be-user/service"
	"google.golang.org/grpc"
	"time"
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

func (s *UserServiceServer) FindOrCreateByStudentId(ctx context.Context,
	request *userv1.FindOrCreateByStudentIdRequest) (*userv1.FindOrCreateByStudentIdResponse, error) {
	u, err := s.svc.FindOrCreateByStudentId(ctx, request.GetStudentId(), request.GetPassword())
	return &userv1.FindOrCreateByStudentIdResponse{
		User: convertToV(u),
	}, err
}

func (s *UserServiceServer) FindByStudentId(ctx context.Context,
	request *userv1.FindByStudentIdRequest) (*userv1.FindByStudentIdResponse, error) {
	u, err := s.svc.FindByStudentId(ctx, request.GetStudentId())
	return &userv1.FindByStudentIdResponse{
		User: convertToV(u),
	}, err
}

func (s *UserServiceServer) FindOrCreateByUserId(ctx context.Context,
	request *userv1.FindOrCreateByUserIdRequest) (*userv1.FindOrCreateByUserIdResponse, error) {
	u, err := s.svc.FindById(ctx, request.GetUserId())
	return &userv1.FindOrCreateByUserIdResponse{
		User: convertToV(u),
	}, err
}

func (s *UserServiceServer) UpdateNonSensitiveInfo(ctx context.Context, request *userv1.UpdateNonSensitiveInfoRequest) (*userv1.UpdateNonSensitiveInfoResponse, error) {
	err := s.svc.UpdateNonSensitiveInfo(ctx, convertToDomain(request.User))
	return &userv1.UpdateNonSensitiveInfoResponse{}, err
}

func (s *UserServiceServer) GetCookie(ctx context.Context, request *userv1.GetCookieRequest) (*userv1.GetCookieResponse, error) {
	u, err := s.svc.GetCookie(ctx, request.GetStudentId())
	if err == service.ErrUserNotFound {
		return &userv1.GetCookieResponse{}, userv1.ErrorUserNotFound("用户不存在: %d", request.GetStudentId())
	}
	return &userv1.GetCookieResponse{Cookie: u}, err
}

func convertToV(user domain.User) *userv1.User {
	return &userv1.User{
		Id:        user.Id,
		StudentId: user.StudentId,
		Password:  user.Password,
		Utime:     user.Utime.UnixMilli(),
		Ctime:     user.Ctime.UnixMilli(),
		New:       user.New,
	}
}

func convertToDomain(user *userv1.User) domain.User {
	return domain.User{
		Id:        user.GetId(),
		StudentId: user.GetStudentId(),
		Password:  user.GetPassword(),
		New:       user.GetCtime() == user.GetUtime(),
		Utime:     time.UnixMilli(user.GetUtime()),
		Ctime:     time.UnixMilli(user.GetCtime()),
	}
}
