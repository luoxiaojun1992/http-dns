package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"github.com/luoxiaojun1992/http-dns/services"
	"github.com/luoxiaojun1992/http-dns/utils"
	"github.com/luoxiaojun1992/DI"
)

var localCache *cache.Cache

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
			if ipListCache, result := localCache.Get("ip:" + QueryObj.Region + ":" + QueryObj.ServiceName); result {
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": gin.H{"ips": ipListCache}})
				return
			}

			ips, err := DI.C.Resolve("ip-service").(*services.IpServiceProto).GetList(QueryObj.Region, QueryObj.ServiceName)

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

	// Basic auth
	authGroup := r.Group("/", gin.BasicAuth(gin.Accounts{os.Getenv("BASIC_AUTH_USER"):os.Getenv("BASIC_AUTH_PWD")}))

	// Ip register
	authGroup.POST("/ip", func(c *gin.Context) {
		var PostForm struct {
			Region      string `form:"region" binding:"required"`
			ServiceName string `form:"service-name" binding:"required"`
			Ip          string `form:"ip" binding:"required"`
			Ttl         string `form:"ttl" binding:"required"`
		}

		err := c.Bind(&PostForm)

		if err == nil {
			ipService := DI.C.Resolve("ip-service").(*services.IpServiceProto)
			_, err := ipService.Add(PostForm.Region, PostForm.ServiceName, PostForm.Ip, PostForm.Ttl)
			if err == nil {
				//Update local cache
				ips, err := ipService.GetList(PostForm.Region, PostForm.ServiceName)
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

	// Ip delete
	authGroup.DELETE("/ip", func(c *gin.Context) {
		var QueryObject struct {
			Region      string `form:"region" binding:"required"`
			ServiceName string `form:"service-name" binding:"required"`
		}

		err := c.Bind(&QueryObject)

		if err == nil {
			ipService := DI.C.Resolve("ip-service").(*services.IpServiceProto)
			_, err := ipService.Delete(QueryObject.Region, QueryObject.ServiceName)
			if err == nil {
				//Update local cache
				ips, err := ipService.GetList(QueryObject.Region, QueryObject.ServiceName)
				if err == nil {
					localCache.Set("ip:"+QueryObject.Region+":"+QueryObject.ServiceName, ips, -1)
				}

				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": err.Error()})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": err.Error()})
		}
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

	var wg sync.WaitGroup

	s := &http.Server{
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	wg.Add(1)
	go func() {
		log.Println(s.Serve(ln))
		wg.Done()
	}()

	//Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	//Graceful Shutdown
	s.Shutdown(nil)

	wg.Wait()
}

func init() {
	//Init env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//Init orm
	utils.InitOrm()

	//Init Local Cache
	localCache = cache.New(1*time.Second, 10*time.Minute)
}

func main() {
	r := setupRouter()

	// Listen and Server in 0.0.0.0:8080
	// r.Run(":9999")
	run(r)
}
