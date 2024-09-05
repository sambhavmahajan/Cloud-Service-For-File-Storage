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

func userExists(uname string) bool {
	mu.RLock()
	_, exists := userData[uname]
	mu.RUnlock()
	return exists
}

func makeNewUser(uname string, pass string) {
	mu.Lock()
	userData[uname] = pass
	mu.Unlock()
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	r.POST("/register", func(c *gin.Context) {
		uname := c.PostForm("username")
		pass := c.PostForm("password")
		if userExists(uname) {
			c.String(http.StatusConflict, "Username already exists!")
			return
		}
		makeNewUser(uname, pass)
		c.String(http.StatusOK, "User Created!")
	})

	r.POST("/login", func(c *gin.Context) {
		uname := c.PostForm("username")
		pass := c.PostForm("password")
		if userExists(uname) {
			mu.RLock()
			if userData[uname] == pass{
				c.String(http.StatusOK, "Login Successful.")
				c.SetCookie("username", uname, 3600, "/", "localhost", false, true)
				c.SetCookie("password", pass, 3600, "/", "localhost", false, true)
			}else {
				c.String(http.StatusConflict, "Bad Login.")
			}
			mu.RUnlock()
		} else {
			c.String(http.StatusUnauthorized, "Bad Credentials.")
		}
	})

	r.Run()
}
