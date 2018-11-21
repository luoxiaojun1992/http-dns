package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"net"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/luoxiaojun1992/http-dns/models"
)

var engine *xorm.Engine

var ipLists map[string][]map[string]string

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Ip resolve
	r.GET("/ips", func(c *gin.Context) {
		var QueryObj struct{
			Region string `form:"region" binding:"required"`
			ServiceName string `form:"service-name" binding:"required"`
		}

		err := c.BindQuery(&QueryObj)

		if err == nil {
			//todo fetch from db & local cache(ttl)
			ips, ok := ipLists[QueryObj.Region+":"+QueryObj.ServiceName]
			if ok {
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": gin.H{"ips": ips}})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "service not found"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"code":1, "msg": err.Error()})
		}
	})

	// Ip register
	r.POST("/ip", func(c *gin.Context) {
		//todo param validation

		serviceID := c.PostForm("region") + ":" + c.PostForm("service-name")
		ips, ok := ipLists[serviceID]
		if ok {
			ipLists[serviceID] = append(ips, map[string]string{"ip": c.PostForm("ip"), "ttl": c.PostForm("ttl")})
		} else {
			ipLists[serviceID] = []map[string]string{{"ip": c.PostForm("ip"), "ttl": "600"}}
		}

		c.JSON(http.StatusOK, gin.H{"code":0, "msg":"ok", "data":gin.H{"ips":ipLists[serviceID]}})
	})

	return r
}

func run(r *gin.Engine) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening:" + port)

	http.Serve(ln, r)

	//todo glance shutdown
}

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:0600120597Abc@/http_dns?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	//Sync Tables
	engine.Sync2(new(models.IpList))
}

func main() {
	ipLists = make(map[string][]map[string]string)
	ipLists["sh:user-service"] = []map[string]string{{"ip": "192.168.0.1", "ttl": "600"}, {"ip": "192.168.0.2", "ttl": "600"}}

	r := setupRouter()

	// Listen and Server in 0.0.0.0:8080
	// r.Run(":9999")
	run(r)
}
