package svc

import (
	"fmt"
	"github.com/copo888/transaction_service/rpc/internal/config"
	"github.com/go-redis/redis/v8"
	"github.com/neccoys/go-driver/mysqlx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strings"
)

type ServiceContext struct {
	Config config.Config
	RedisClient *redis.Client
	MyDB        *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {

	redisCache := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    c.RedisCache.RedisMasterName,
		SentinelAddrs: strings.Split(c.RedisCache.RedisSentinelNode, ";"),
		DB:            c.RedisCache.RedisDB,
	})

	log.Println(c.Mysql)
	myDb, err := mysqlx.New(c.Mysql.Host, fmt.Sprintf("%d", c.Mysql.Port), c.Mysql.UserName, c.Mysql.Password, c.Mysql.DBName).
		//SetCharset("utf8mb4").
		SetLoc("UTC").
		SetLogger(logger.Default.LogMode(logger.Info)).
		Connect(mysqlx.Pool(10, 100, 180))

	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config: c,
		MyDB: myDb,
		RedisClient: redisCache,
	}
}
