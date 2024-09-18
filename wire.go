//go:build wireinject

package main

import (
	"github.com/asynccnu/be-user/grpc"
	"github.com/asynccnu/be-user/ioc"
	"github.com/asynccnu/be-user/pkg/grpcx"
	"github.com/asynccnu/be-user/repository"
	"github.com/asynccnu/be-user/repository/cache"
	"github.com/asynccnu/be-user/repository/dao"
	"github.com/asynccnu/be-user/service"
	"github.com/google/wire"
)

func InitGRPCServer() grpcx.Server {
	wire.Build(
		ioc.InitGRPCxKratosServer,
		grpc.NewUserServiceServer,
		service.NewUserService,
		repository.NewCachedUserRepository,
		dao.NewGORMUserDAO,
		cache.NewRedisUserCache,
		// 第三方
		ioc.InitCCNUClient,
		ioc.InitEtcdClient,
		ioc.InitRedis,
		ioc.InitDB,
		ioc.InitLogger,
	)
	return grpcx.Server(nil)
}
