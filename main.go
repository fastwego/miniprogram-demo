package main

import (
	"context"
	"fmt"
	"github.com/fastwego/miniprogram"
	"github.com/fastwego/miniprogram/apis/datacube"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

}
func main() {

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	app := miniprogram.New(miniprogram.MiniprogramConfig{
		Appid:  viper.GetString("APPID"),
		Secret: viper.GetString("SECRET"),
	})

	// 接口演示
	router.GET("/api/weixin/miniprogram", func(c *gin.Context) {
		var payload = []byte(`{
  "begin_date" : "20170313",
  "end_date" : "20170313"
}`)
		resp, err := datacube.GetDailyRetain(app, payload)
		fmt.Println(string(resp),err)

		c.Writer.Write(resp)
	})

	svr := &http.Server{
		Addr:    viper.GetString("LISTEN"),
		Handler: router,
	}

	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeout := time.Duration(5) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
