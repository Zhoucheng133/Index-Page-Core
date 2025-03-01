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
	r.Use(cors.Default())

	r.GET("/api/list", utils.List)
	r.POST("/api/add", utils.AddItem)
	r.DELETE("/api/del/:id", utils.DeleteItem)

	r.GET("/api/init", utils.Init)

	r.Run(":8080")
}
