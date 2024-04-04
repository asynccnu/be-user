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
	LoginByCCNU(ctx context.Context, studentId string, password string) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, uid int64) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	ccnu CCNUService
}

func (s *userService) FindById(ctx context.Context, uid int64) (domain.User, error) {
	return s.repo.FindById(ctx, uid)
}

func (s *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return s.repo.UpdateSensitiveInfo(ctx, user)
}

func (s *userService) LoginByCCNU(ctx context.Context, studentId string, password string) (domain.User, error) {
	// 模拟登录
	ok, err := s.ccnu.Login(ctx, studentId, password)
	if err != nil {
		return domain.User{}, err
	}
	if !ok {
		return domain.User{}, ErrInvalidStudentIdOrPassword
	}

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

func NewUserService(repo repository.UserRepository, ccnu CCNUService) UserService {
	return &userService{repo: repo, ccnu: ccnu}
}
