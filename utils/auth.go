package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jaevor/go-nanoid"
)

var secretKey = []byte("index")

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Salt     string `json:"salt"`
}

func userExist() bool {
	rows, err := db.Query("SELECT id, name, password FROM users")
	if err != nil {
		return false
	}
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Password,
		); err != nil {
			return false
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return false
	} else {
		return true
	}
}

// 初始化判断是否有用户信息
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
			c.JSON(200, gin.H{"error": fmt.Sprint("数据处理错误", err)})
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

// 保存密码
func savePassword(password string, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

// 注册
func Register(c *gin.Context) {

	if userExist() {
		c.JSON(200, gin.H{"ok": false, "msg": "用户已存在"})
		return
	}

	var newUser User
	if err := c.ShouldBind(&newUser); err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "请求数据格式不正确"})
		return
	}
	query := `INSERT INTO users (id, name, password, salt) VALUES (?, ?, ?, ?)`
	id, _ := nanoid.Standard(21)
	salt, _ := nanoid.Standard(10)
	saltString := salt()
	_, err := db.Exec(query, id(), newUser.Name, savePassword(newUser.Password, saltString), saltString)
	if err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "注册失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"ok": true, "msg": "注册成功"})
}

// 生成token
func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 365).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// 登录验证
func LoginCheck(username string, password string) bool {
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

// token验证
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		excludedPaths := []string{"/api/login", "/api/register", "/api/init"}
		path := c.Request.URL.Path
		for _, excludedPath := range excludedPaths {
			if strings.HasPrefix(path, excludedPath) {
				c.Next()
				return
			}
		}

		tokenString := c.GetHeader("auth")
		if tokenString == "" {
			c.JSON(200, gin.H{"ok": false, "msg": "Missing Authorization Header"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(200, gin.H{"ok": false, "msg": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 登录
func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "请求数据格式不正确"})
		return
	}
	if LoginCheck(user.Name, user.Password) {
		token, err := GenerateToken(user.Name)
		if err != nil {
			c.JSON(200, gin.H{"ok": false, "msg": err})
		}
		c.JSON(200, gin.H{"ok": true, "msg": token})
	} else {
		c.JSON(200, gin.H{"ok": false, "msg": "身份验证失败"})
	}
}
