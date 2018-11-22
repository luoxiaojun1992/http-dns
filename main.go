package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/luoxiaojun1992/http-dns/models"
	"github.com/patrickmn/go-cache"
	"time"
)

var orm *xorm.Engine
var localCache *cache.Cache

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
		var QueryObj struct {
			Region      string `form:"region" binding:"required"`
			ServiceName string `form:"service-name" binding:"required"`
		}

		err := c.BindQuery(&QueryObj)

		if err == nil {
			ips := make([]models.IpList, 0, 10)

			if ipListCache, result := localCache.Get("ip:"+QueryObj.Region+":"+QueryObj.ServiceName); result {
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": gin.H{"ips": ipListCache}})
				return
			}

			err := orm.Where("region = ? AND service_name = ?", QueryObj.Region, QueryObj.ServiceName).
				Limit(10).
				Select("ip, ttl").
				Find(&ips)
			if err == nil {
				localCache.Set("ip:"+QueryObj.Region+":"+QueryObj.ServiceName, ips, -1)
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": gin.H{"ips": ips}})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": err.Error()})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": err.Error()})
		}
	})

	// Ip register
	r.POST("/ip", func(c *gin.Context) {
		var PostForm struct {
			Region      string `form:"region" binding:"required"`
			ServiceName string `form:"service-name" binding:"required"`
			Ip          string `form:"ip" binding:"required"`
			Ttl         string `form:"ttl" binding:"required"`
		}

		err := c.Bind(&PostForm)

		if err == nil {
			_, err := orm.Insert(models.IpList{
				Region:      PostForm.Region,
				ServiceName: PostForm.ServiceName,
				Ip:          PostForm.Ip,
				Ttl:         PostForm.Ttl,
			})

			if err == nil {
				//Update local cache
				ips := make([]models.IpList, 0, 10)

				err := orm.Where("region = ? AND service_name = ?", PostForm.Region, PostForm.ServiceName).
					Limit(10).
					Select("ip, ttl").
					Find(&ips)
				if err == nil {
					localCache.Set("ip:"+PostForm.Region+":"+PostForm.ServiceName, ips, -1)
				}

				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": err.Error()})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": err.Error()})
		}
	})

	//todo delete

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
	//Init ORM
	var err error
	//todo env or config
	orm, err = xorm.NewEngine("mysql", "root:0600120597$Abc@/http_dns?charset=utf8mb4")
	if err != nil {
		log.Fatal(err)
	}

	//Sync Tables
	err = orm.Sync2(new(models.IpList))
	if err != nil {
		log.Fatal(err)
	}

	//Init Local Cache
	localCache = cache.New(1*time.Second, 10*time.Minute)
}

func main() {
	ipLists = make(map[string][]map[string]string)
	ipLists["sh:user-service"] = []map[string]string{{"ip": "192.168.0.1", "ttl": "600"}, {"ip": "192.168.0.2", "ttl": "600"}}

	r := setupRouter()

	// Listen and Server in 0.0.0.0:8080
	// r.Run(":9999")
	run(r)
}
