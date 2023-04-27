package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
    "log"

	{{.importPackages}}
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
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
