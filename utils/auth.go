package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type UserIn struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Init(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, password FROM users")
	if err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "数据处理错误"})
	}
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Password,
		); err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprint("数据处理错误", err)})
			return
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		c.JSON(200, gin.H{"ok": true, "msg": true})
	} else {
		c.JSON(200, gin.H{"ok": true, "msg": false})
	}
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result []byte
	for i := 0; i < length; i++ {
		randomIndex := rng.Intn(len(charset))
		result = append(result, charset[randomIndex])
	}
	return string(result)
}

func savePassword(password string, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

func Register(c *gin.Context) {
	var newUser UserIn
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(200, gin.H{"ok": false, "data": "请求数据格式不正确"})
		return
	}
	query := `INSERT INTO users (name, password, salt) VALUES (?, ?, ?)`
	salt := generateRandomString(6)
	_, err := db.Exec(query, newUser.Name, savePassword(newUser.Password, salt), salt)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "data": "注册失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"ok": true, "data": "注册成功"})
}

func AuthCheck(username string, password string) bool {
	rows, err := db.Query("SELECT id, name, password, salt FROM users WHERE name= ? ", username)
	if err != nil {
		return false
	}
	var userData User
	if rows.Next() {
		if err := rows.Scan(&userData.ID, &userData.Name, &userData.Password, &userData.Salt); err != nil {
			defer rows.Close()
			return false
		} else if savePassword(password, userData.Salt) == userData.Password {
			defer rows.Close()
			return true
		} else {
			defer rows.Close()
			return false
		}
	} else {
		defer rows.Close()
		return false
	}
}

func Login(c *gin.Context) {
	var user UserIn
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{"ok": false, "data": "请求数据格式不正确"})
		return
	}
	if AuthCheck(user.Name, user.Password) {
		c.JSON(200, gin.H{"ok": true, "msg": ""})
	} else {
		c.JSON(200, gin.H{"ok": false, "msg": "身份验证失败"})
	}
}
