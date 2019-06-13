package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
	"strconv"
	"strings"
)

var (
	redisdb *redis.ClusterClient
	cachedb *cache.Cache
)

type redisConfiguration struct {
	Hosts    string
	Password string
	DataNum  int8
}

func initRedisConnect() error {
	var config redisConfiguration
	if err := envconfig.Process("redis", &config); err != nil {
		logs.Error(err, "Startup")
		return err
	}
	if isBlank := CheckStringIsBlank(config.Hosts); isBlank {
		config.Hosts = beego.AppConfig.String("redisHost")
		dataNum, er := strconv.ParseInt(beego.AppConfig.String("redisNum"), 10, 8)
		if er != nil {
			logs.Warning(er, "redis database number parse error, user default 0")
			dataNum = 0
		}
		config.DataNum = int8(dataNum)
	}
	redisOps := new(redis.ClusterOptions)
	redisOps.Addrs = strings.Split(config.Hosts, ",")
	if isBlank := CheckStringIsBlank(config.Password); isBlank {
		redisOps.Password = config.Password
	}
	redisdb = redis.NewClusterClient(redisOps)
	logs.Debug("init redis connect")
	return nil
}

func GetCacheClient() *redis.ClusterClient {
	return redisdb
}

func shutdownRedisConnect() {
	if redisdb != nil {
		logs.Debug("redis Disconnect")
		if err := redisdb.Close(); err != nil {
			logs.Error("redis close error", err)
		}
	}
}

func initCache() error {
	bm, err := cache.NewCache("redis", `{"conn":"127.0.0.1:6379","key":"collectionName","dbNum":"0","password":""}`)
	if err != nil {
		logs.Error("cache init error", err)
		return err
	} else {
		cachedb = &bm
		return nil
	}
}

func CloseCache() {
	if cachedb != nil {

	}
}
