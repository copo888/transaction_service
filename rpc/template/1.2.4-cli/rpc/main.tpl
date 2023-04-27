package main

import (
	"flag"
	"fmt"
    {{if .consul}}"github.com/neccoys/go-zero-extension/consul"{{end}}
    {{if .check}}"google.golang.org/grpc/health/grpc_health_v1"{{end}}
    "log"

	{{.imports}}

    "github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
    configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")
    envFile    = flag.String("env", "etc/.env", "the env file")
)

func main() {
	flag.Parse()

    if err := godotenv.Load(*envFile); err != nil {
        log.Fatal("Error loading .env file")
    }

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)
	srv := server.New{{.serviceNew}}Server(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		{{.pkg}}.Register{{.service}}Server(grpcServer, srv)
        {{if .check}}grpc_health_v1.RegisterHealthServer(grpcServer, srv){{end}}

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	{{if .consul}}
	// 注册Consul服务
    if err := consul.RegisterService(c.ListenOn, c.Consul); err != nil {
        log.Println("Consul Error:", err)
    }
    {{end}}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
