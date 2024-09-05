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
				c.SetCookie("username", uname, 3600, "/", "localhost", false, true)
				c.SetCookie("password", pass, 3600, "/", "localhost", false, true)
				c.Redirect(http.StatusPermanentRedirect, "/user/"+uname)
			}else {
				c.String(http.StatusConflict, "Bad Login.")
			}
			mu.RUnlock()
		} else {
			c.String(http.StatusUnauthorized, "Bad Credentials.")
		}
	})

	r.GET("/user/:username", func(c *gin.Context){
		uname := c.Param("username")
		cache_name, err := c.Cookie("username")
		if err != nil{
			c.String(http.StatusUnauthorized, "Can't view without login")
			return
		}
		if uname != cache_name{
			c.String(http.StatusUnauthorized, "Can't view without login")
			return
		}
		mu.RLock()
		val, isValid := userData[uname]
		cache_pass, err1 := c.Cookie("password")
		if err1 != nil || !isValid || val != cache_pass {
			c.String(http.StatusUnauthorized, "Can't view without login")
			return
		}
		c.String(200, "Login as " + uname)
		mu.RUnlock()
	})

	r.Run()
}
