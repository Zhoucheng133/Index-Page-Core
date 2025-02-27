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

	r.Run(":8080")
}
