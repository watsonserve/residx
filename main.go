package main

import (
	"fmt"
	"os"

	"github.com/watsonserve/goengine"
	"github.com/watsonserve/scaner/actions"
	"github.com/watsonserve/scaner/utils"
)

func main() {
	// load configure
	cfg, err := utils.GetOption()
	if nil != err {
		fmt.Fprint(os.Stderr, err.Error())
	}
	lis := cfg["listen"][0]
	redisAddr := cfg["redis"][0]
	dbCfg := &goengine.DbConf{
		User:   cfg["db_user"][0],
		Passwd: cfg["db_passwd"][0],
		Host:   cfg["db_host"][0],
		Port:   cfg["db_port"][0],
		Name:   cfg["db_name"][0],
	}

	// connect to mongodb
	db := goengine.ConnMongo(dbCfg)
	if nil == db {
		fmt.Fprint(os.Stderr, "Connect to mongodb failed")
		return
	}

	// redis client and session manager
	var sm goengine.SessionManager = nil
	if "" != redisAddr {
		redis := goengine.NewRedisStore(redisAddr, "", 0)
		if nil == redis {
			fmt.Fprint(os.Stderr, "Connect to redis failed")
			return
		}
		sm = goengine.InitSessionManager(redis, "_rs", "sess:", "rsch:", "")
	}

	ac := actions.New(db, cfg["root"][0])
	router := goengine.InitHttpRoute()
	ac.Bind(router)

	// set up http server
	app := goengine.New(sm)
	app.UseRouter(router)
	if "../" == lis[0:3] || "./" == lis[0:2] || '/' == lis[0] {
		app.ListenUnix(lis)
	} else {
		app.ListenTCP(lis)
	}
}
