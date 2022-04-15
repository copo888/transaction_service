package main

import (
	"flag"
	"fmt"
	"github.com/neccoys/go-zero-extension/consul"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"

	"github.com/copo888/transaction_service/rpc/internal/config"
	"github.com/copo888/transaction_service/rpc/internal/server"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/transaction"

	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configFile = flag.String("f", "rpc/etc/transaction.yaml", "the config file")
	envFile    = flag.String("env", "rpc/etc/.env", "the env file")
)

func main() {
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Fatal("Error loading .env file")
	}

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)
	srv := server.NewTransactionServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		transaction.RegisterTransactionServer(grpcServer, srv)
		grpc_health_v1.RegisterHealthServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 注册Consul服务
	if err := consul.RegisterService(c.ListenOn, c.Consul); err != nil {
		log.Println("Consul Error:", err)
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
