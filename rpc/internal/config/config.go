package config

import (
	"github.com/neccoys/go-zero-extension/consul"
	_ "github.com/zeromicro/go-zero/core/service"
	_ "github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Consul consul.Conf
	Mysql  struct {
		Host       string
		Port       int
		DBName     string
		UserName   string
		Password   string
		DebugLevel string
	}
	RedisCache struct {
		RedisSentinelNode string
		RedisMasterName   string
		RedisDB           int
	}
	Version string
}
