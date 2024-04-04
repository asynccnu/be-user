package repository

import (
	"context"
	"github.com/MuxiKeStack/be-user/domain"
	"github.com/MuxiKeStack/be-user/repository/cache"
	"github.com/MuxiKeStack/be-user/repository/dao"
	"time"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateUser
	ErrUserNotFind   = dao.ErrRecordNotFind
)

type UserRepository interface {
	FindByStudentId(ctx context.Context, studentId string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	UpdateSensitiveInfo(ctx context.Context, user domain.User) error
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func (repo *CachedUserRepository) UpdateSensitiveInfo(ctx context.Context, user domain.User) error {
	err := repo.dao.UpdateSensitiveInfoById(ctx, repo.toEntity(user))
	if err != nil {
		return err
	}
	return nil
	// TODO 认为用户短期可能重新查看修改的内容，进行预缓存
	//repo.cache.Set(ctx)
}

func (repo *CachedUserRepository) FindByStudentId(ctx context.Context, studentId string) (domain.User, error) {
	u, err := repo.dao.FindByStudentId(ctx, studentId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:        u.Id,
		StudentId: u.Sid,
		Avatar:    u.Avatar,
		Nickname:  u.Nickname,
		New:       u.Utime == u.Ctime, // 更新时间为创建时间说明是未更新过信息的新用户
		Utime:     time.UnixMilli(u.Utime),
		Ctime:     time.UnixMilli(u.Ctime),
	}
}
func (repo *CachedUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Sid:      u.StudentId,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
	}
}

func NewUserRepository(dao dao.UserDAO) UserRepository {
	return &CachedUserRepository{
		dao: dao,
	}
}
