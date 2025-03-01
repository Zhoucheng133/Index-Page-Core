package utils

import "github.com/gin-gonic/gin"

type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Init(c *gin.Context) {
	
}
