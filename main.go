package main

import (
	"index_page_core/routers"
	"index_page_core/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitSql()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(cors.Default())

	r.GET("/api/list", routers.List)
	r.POST("/api/add", routers.AddItem)
	r.DELETE("/api/del/:id", routers.DeleteItem)

	r.Run(":8080")
}
