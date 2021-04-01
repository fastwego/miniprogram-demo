// Copyright 2020 FastWeGo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fastwego/miniprogram"
	"github.com/fastwego/miniprogram/apis/datacube"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

var App *miniprogram.Miniprogram

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	App = miniprogram.New(miniprogram.Config{
		Appid:  viper.GetString("APPID"),
		Secret: viper.GetString("SECRET"),
	})

}
func main() {

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 接口演示
	router.GET("/api/weixin/miniprogram", func(c *gin.Context) {
		var payload = []byte(`{
		  "begin_date" : "20170313",
		  "end_date" : "20170313"
		}`)
		resp, err := datacube.GetDailyRetain(App, payload)
		fmt.Println(string(resp), err)

		c.Writer.Write(resp)
	})

	router.GET("/api/img_sec_check", ImgSecCheck)

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
