package config

import (
    "github.com/zeromicro/go-zero/zrpc"
    {{if .consul}}"github.com/neccoys/go-zero-extension/consul"{{end}}
)

type Config struct {
	zrpc.RpcServerConf
    {{if .consul}}Consul consul.Conf{{end}}
}
