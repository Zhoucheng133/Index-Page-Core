package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("index")

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

// 生成salt
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

// 保存密码
func savePassword(password string, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

// 注册
func Register(c *gin.Context) {
	var newUser UserIn
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "请求数据格式不正确"})
		return
	}
	query := `INSERT INTO users (name, password, salt) VALUES (?, ?, ?)`
	salt := generateRandomString(6)
	_, err := db.Exec(query, newUser.Name, savePassword(newUser.Password, salt), salt)
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
	var user UserIn
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "请求数据格式不正确"})
		return
	}
	if LoginCheck(user.Name, user.Password) {
		token, err := GenerateToken(user.Name)
		fmt.Println(token)
		if err != nil {
			c.JSON(200, gin.H{"ok": false, "msg": err})
		}
		c.JSON(200, gin.H{"ok": true, "msg": token})
	} else {
		c.JSON(200, gin.H{"ok": false, "msg": "身份验证失败"})
	}
}
