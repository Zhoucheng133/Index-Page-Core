package utils

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jaevor/go-nanoid"
	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	ID    string `json:"id"`
	Name  string `json:"name" binding:"required"`
	Port  string `json:"port" binding:"required"`
	WebUI string `json:"webui" binding:"required"`
	Tip   string `json:"tip" binding:"omitempty"`
}

func List(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, port, webui, tip FROM pages")
	if err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": err})
		return
	}
	var pages []Page
	for rows.Next() {
		var p Page
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Port,
			&p.WebUI,
			&p.Tip,
		); err != nil {
			log.Printf("数据解析失败: %v", err)
			c.JSON(200, gin.H{"error": fmt.Sprint("数据处理错误", err)})
			return
		}
		pages = append(pages, p)
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": err})
		return
	}
	if pages != nil {
		c.JSON(200, gin.H{"ok": true, "msg": pages})
	} else {
		c.JSON(200, gin.H{"ok": true, "msg": []Page{}})
	}

}

func AddItem(c *gin.Context) {
	var newPage Page
	if err := c.ShouldBindJSON(&newPage); err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "请求数据格式不正确"})
		return
	}

	id, _ := nanoid.Standard(21)
	idString := id()
	query := `INSERT INTO pages (id, name, port, webui, tip) VALUES(?, ?, ?, ?, ?)`
	_, err := db.Exec(query, idString, newPage.Name, newPage.Port, newPage.WebUI, newPage.Tip)
	if err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "插入数据失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"ok": true, "msg": "数据插入成功", "id": idString})
}

func DeleteItem(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(200, gin.H{"ok": false, "msg": "缺少 ID 参数"})
		return
	}

	query := `DELETE FROM pages WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "删除数据失败", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true, "msg": "数据删除成功"})
}

func EditItem(c *gin.Context) {
	var newPage Page
	if err := c.ShouldBindJSON(&newPage); err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "请求数据格式不正确"})
		return
	}
	query := `UPDATE pages SET name = ?, port = ?, webui = ?, tip = ? WHERE id = ?`
	_, err := db.Exec(query, newPage.Name, newPage.Port, newPage.WebUI, newPage.Tip, newPage.ID)
	if err != nil {
		c.JSON(200, gin.H{"ok": false, "msg": "更新数据失败", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true, "msg": "更新数据成功"})
}
