package utils

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetIpv4(c *gin.Context) {
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
