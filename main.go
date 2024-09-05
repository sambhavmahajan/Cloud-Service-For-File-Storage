package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var (
	userData = make(map[string]string)

	mu sync.RWMutex
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	r.POST("/register", func(c *gin.Context) {
		uname := c.PostForm("username")
		pass := c.PostForm("password")
		mu.RLock()
		_, exists := userData[uname]
		mu.RUnlock()
		if exists {
			c.String(http.StatusConflict, "Username already exists!")
			return
		}
		mu.Lock()
		userData[uname] = pass
		mu.Unlock()
		c.String(http.StatusOK, "User Created!")
	})
	r.Run()
}
