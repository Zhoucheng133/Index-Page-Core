package routers

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
	db, err := sql.Open("sqlite3", "db/pages.db")
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
	defer db.Close()
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
	defer rows.Close()
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
	var newPage Page
	if err := c.ShouldBindJSON(&newPage); err != nil {
		c.JSON(400, gin.H{"ok": false, "data": "请求数据格式不正确"})
		return
	}
	db, err := sql.Open("sqlite3", "db/pages.db")
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "data": "数据库连接失败"})
		return
	}
	defer db.Close()
	query := `INSERT INTO pages (name, port, webui, tip) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, newPage.Name, newPage.Port, newPage.WebUI, newPage.Tip)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "data": "插入数据失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"ok": true, "data": "数据插入成功"})
}
