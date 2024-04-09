package service

import (
	"context"
	"errors"
	"github.com/MuxiKeStack/be-user/domain"
	"github.com/MuxiKeStack/be-user/repository"
)

var (
	ErrInvalidStudentIdOrPassword = errors.New("学号或密码不对")
)

type UserService interface {
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindOrCreateByStudentId(ctx context.Context, studentId string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func (s *userService) FindOrCreateByStudentId(ctx context.Context, studentId string) (domain.User, error) {
	u, err := s.repo.FindByStudentId(ctx, studentId)
	if err == nil {
		return u, nil
	}
	// 系统异常，返回错误
	if err != repository.ErrUserNotFind {
		return domain.User{}, err
	}
	// 用户不存在，首次登录，创建用户
	err = s.repo.Create(ctx, domain.User{StudentId: studentId})
	// 并发场景下，如果错误为非duplicate错误，则为系统异常
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	// 如果后续分库分表，这里必须从主库查询
	return s.repo.FindByStudentId(ctx, studentId)
}

func (s *userService) FindById(ctx context.Context, uid int64) (domain.User, error) {
	return s.repo.FindById(ctx, uid)
}

func (s *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return s.repo.UpdateSensitiveInfo(ctx, user)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
