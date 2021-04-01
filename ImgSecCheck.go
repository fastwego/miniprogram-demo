package main

import (
	"fmt"
	"github.com/fastwego/miniprogram/apis/security"
	"github.com/gin-gonic/gin"
)

func ImgSecCheck(c *gin.Context) {

	resp, err := security.ImgSecCheck(App,"data/hi.jpg")

	fmt.Println(string(resp), err)

	c.Writer.Write(resp)
}
