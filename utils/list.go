package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name"`
	Port  string `json:"port"`
	WebUI int    `json:"webui"`
	Tip   string `json:"tip"`
}

func List(c *gin.Context) {
	username := c.GetHeader("username")
	password := c.GetHeader("password")
	if len(username) == 0 || len(password) == 0 {
		c.JSON(401, gin.H{"ok": false, "msg": "缺少请求头"})
		return
	} else if !AuthCheck(username, password) {
		c.JSON(401, gin.H{"ok": false, "msg": "身份验证失败"})
		return
	}
	rows, err := db.Query("SELECT id, name, port, webui, tip FROM pages")
	if err != nil {
		c.JSON(
			400,
			gin.H{
				"ok":   false,
				"data": err,
			},
		)
		return
	}
	var pages []Page
	for rows.Next() {
		var p Page
		var tip sql.NullString
		// 扫描顺序与SELECT字段顺序严格对应
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Port,
			&p.WebUI,
			&tip,
		); err != nil {
			log.Printf("数据解析失败: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprint("数据处理错误", err)})
			return
		}
		if tip.Valid {
			p.Tip = tip.String
		} else {
			p.Tip = ""
		}
		pages = append(pages, p)
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		c.JSON(
			400,
			gin.H{
				"ok":   false,
				"data": err,
			},
		)
		return
	}
	c.JSON(
		200,
		gin.H{
			"ok":   true,
			"data": pages,
		},
	)
}

func AddItem(c *gin.Context) {
	username := c.GetHeader("username")
	password := c.GetHeader("password")
	if len(username) == 0 || len(password) == 0 {
		c.JSON(401, gin.H{"ok": false, "msg": "缺少请求头"})
		return
	} else if !AuthCheck(username, password) {
		c.JSON(401, gin.H{"ok": false, "msg": "身份验证失败"})
		return
	}
	var newPage Page
	if err := c.ShouldBindJSON(&newPage); err != nil {
		c.JSON(400, gin.H{"ok": false, "data": "请求数据格式不正确"})
		return
	}
	query := `INSERT INTO pages (name, port, webui, tip) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, newPage.Name, newPage.Port, newPage.WebUI, newPage.Tip)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "data": "插入数据失败", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true, "data": "数据插入成功"})
}

func DeleteItem(c *gin.Context) {
	username := c.GetHeader("username")
	password := c.GetHeader("password")
	if len(username) == 0 || len(password) == 0 {
		c.JSON(401, gin.H{"ok": false, "msg": "缺少请求头"})
		return
	} else if !AuthCheck(username, password) {
		c.JSON(401, gin.H{"ok": false, "msg": "身份验证失败"})
		return
	}
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"ok": false, "data": "缺少 ID 参数"})
		return
	}

	query := `DELETE FROM pages WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "data": "删除数据失败", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true, "data": "数据删除成功"})
}
