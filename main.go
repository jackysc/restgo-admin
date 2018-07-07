package main

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/go-xorm/xorm"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"

	"os"
	"strconv"

	"github.com/jackysc/restgo-admin/controller"
	"github.com/jackysc/restgo-admin/entity"
	"github.com/jackysc/restgo-admin/restgo"
)

func registerRouter(router *gin.Engine) {
	new(controller.PageController).Router(router)
	new(controller.UserController).Router(router)
	new(controller.OpenController).Router(router)
	new(controller.RoleController).Router(router)
	new(controller.ResController).Router(router)
	new(controller.ConfigController).Router(router)

}

func init() {
	gob.Register(entity.User{})
}

func main() {

	cfg := new(restgo.Config)
	cfg.Parse("config/app.properties")
	fmt.Println("[ok] load config ")
	restgo.SetCfg(cfg)

	restgo.Configuration(cfg.Logger["filepath"])

	gin.SetMode(cfg.App["mode"])

	for k, ds := range cfg.Datasource {
		e, err := xorm.NewEngine(ds["driveName"], ds["dataSourceName"])
		if err != nil {
			fmt.Println("data source init error", err.Error())
			return
		}
		fmt.Printf("initt data source %s", ds["dataSourceName"])
		e.ShowSQL(ds["showSql"] == "true")
		n, _ := strconv.Atoi(ds["maxIdle"])
		e.SetMaxIdleConns(n)
		n, _ = strconv.Atoi(ds["maxOpen"])
		e.SetMaxOpenConns(n)
		//判断init文件是否存在
		_, err = os.Stat("inited")
		//如果不存在

		if !(err == nil || !os.IsNotExist(err)) {
			fmt.Println("init table and passwd")
			//创建表
			err = e.Sync2(new(entity.User), new(entity.Config), new(entity.RefRoleRes), new(entity.Resource), new(entity.Role))
			if err != nil {
				fmt.Println("data source init error", err.Error())
				return
			}
			//初始化sql语句
			initsql := "INSERT INTO `user` VALUES (1,'admin','18600000000','d060812a3a1af12643a74a4d3b6d492d','admin@qq.com',NULL,'2018-02-23 11:32:32','winlion',0,'admin@qq.com',1)"
			e.Exec(initsql)
			//创建一个文件
			os.Create("inited")

		}

		restgo.SetEngin(k, e)
	}
	fmt.Println("[ok] init datasource")
	router := gin.Default()

	for k, v := range cfg.Static {
		router.Static(k, v)
	}
	for k, v := range cfg.StaticFile {
		router.StaticFile(k, v)
	}

	router.SetFuncMap(restgo.GetFuncMap())
	router.NoRoute(restgo.NoRoute)
	router.NoMethod(restgo.NoMethod)

	router.LoadHTMLGlob(cfg.View["path"] + "/**/*")
	router.Delims(cfg.View["deliml"], cfg.View["delimr"])
	store, err := redis.NewStore(10, "tcp", fmt.Sprintf("%s:%s", cfg.Redis["host"], cfg.Redis["port"]), cfg.Redis["passwd"], []byte(cfg.Session["name"]))
	if err != nil {
		fmt.Println(err.Error())
	}
	router.Use(sessions.Sessions(cfg.Session["name"], store))
	router.Use(restgo.Auth())
	registerRouter(router)

	err = http.ListenAndServe(cfg.App["addr"]+":"+cfg.App["port"], router)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("[ok] app run", cfg.App["addr"]+":"+cfg.App["port"])
	}
}
