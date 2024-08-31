package repository

import (
	"context"
	"github.com/asynccnu/be-user/domain"
	"github.com/asynccnu/be-user/pkg/logger"
	"github.com/asynccnu/be-user/repository/cache"
	"github.com/asynccnu/be-user/repository/dao"
	"time"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateUser
	ErrUserNotFind   = dao.ErrRecordNotFind
)

type UserRepository interface {
	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindByStudentId(ctx context.Context, studentId string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	UpdateSensitiveInfo(ctx context.Context, user domain.User) error
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
	l     logger.Logger
}

func NewCachedUserRepository(dao dao.UserDAO, cache cache.UserCache, l logger.Logger) UserRepository {
	return &CachedUserRepository{dao: dao, cache: cache, l: l}
}

func (repo *CachedUserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	res, err := repo.cache.Get(ctx, uid)
	if err == nil {
		return res, nil
	}
	if err != cache.ErrKeyNotExists {
		// redis崩溃或者网络错误，用户量不大，MySQL撑得住，所以不降级处理
		repo.l.Error("访问Redis失败，查询用户缓存", logger.Error(err), logger.Int64("uid", uid))
	}
	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	res = repo.toDomain(u) // 这里不适合回写
	return res, nil
}

func (repo *CachedUserRepository) UpdateSensitiveInfo(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateSensitiveInfoById(ctx, repo.toEntity(user))
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
