package main

import (
	"index_page_core/routers"
	"index_page_core/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitSql()
	r := gin.New()

	r.GET("/api/list", routers.List)
	r.POST("/api/add", routers.AddItem)
	r.DELETE("/api/del/:id", routers.DeleteItem)

	r.Run(":8080")
}
