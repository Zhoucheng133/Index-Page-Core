package utils

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetIpv4(c *gin.Context) {
	username := c.GetHeader("name")
	password := c.GetHeader("password")
	if len(username) == 0 || len(password) == 0 {
		c.JSON(200, gin.H{"ok": false, "msg": "缺少请求头"})
		return
	} else if !AuthCheck(username, password) {
		c.JSON(200, gin.H{"ok": false, "msg": "身份验证失败"})
		return
	}
	resp, err := http.Get("https://4.ipw.cn")
	if err != nil {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": err,
		})
		return
	}
	body, _ := io.ReadAll(resp.Body)
	c.JSON(200, gin.H{
		"ok":  true,
		"msg": string(body),
	})
}

func GetIpv6(c *gin.Context) {
	username := c.GetHeader("name")
	password := c.GetHeader("password")
	if len(username) == 0 || len(password) == 0 {
		c.JSON(200, gin.H{"ok": false, "msg": "缺少请求头"})
		return
	} else if !AuthCheck(username, password) {
		c.JSON(200, gin.H{"ok": false, "msg": "身份验证失败"})
		return
	}
	resp, err := http.Get("https://6.ipw.cn")
	if err != nil {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": err,
		})
		return
	}
	body, _ := io.ReadAll(resp.Body)
	c.JSON(200, gin.H{
		"ok":  true,
		"msg": string(body),
	})
}
