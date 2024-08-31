// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/asynccnu/be-user/grpc"
	"github.com/asynccnu/be-user/ioc"
	"github.com/asynccnu/be-user/pkg/grpcx"
	"github.com/asynccnu/be-user/repository"
	"github.com/asynccnu/be-user/repository/cache"
	"github.com/asynccnu/be-user/repository/dao"
	"github.com/asynccnu/be-user/service"
)

// Injectors from wire.go:

func InitGRPCServer() grpcx.Server {
	logger := ioc.InitLogger()
	db := ioc.InitDB(logger)
	userDAO := dao.NewGORMUserDAO(db)
	cmdable := ioc.InitRedis()
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache, logger)
	userService := service.NewUserService(userRepository)
	userServiceServer := grpc.NewUserServiceServer(userService)
	client := ioc.InitEtcdClient()
	server := ioc.InitGRPCxKratosServer(userServiceServer, client, logger)
	return server
}
