package main

import (
	"index_page_core/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitSql()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "name", "password")
	r.Use(cors.New(config))

	r.GET("/api/list", utils.List)
	r.POST("/api/add", utils.AddItem)
	r.DELETE("/api/del/:id", utils.DeleteItem)
	r.POST("/api/edit", utils.EditItem)

	r.GET("/api/init", utils.Init)
	r.POST("/api/register", utils.Register)
	r.POST("/api/login", utils.Login)

	r.GET("/api/ipv4", utils.GetIpv4)
	r.GET("/api/ipv6", utils.GetIpv6)

	r.Run(":8080")
}
