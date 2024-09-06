package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

var (
	userData = make(map[string]string)
	usernameToLinks = make(map[string][]string)
	mu sync.RWMutex
)

func userExists(uname string) bool {
	mu.RLock()
	_, exists := userData[uname]
	mu.RUnlock()
	return exists
}

func makeNewUser(uname string, upass string) {
	mu.Lock()
	userData[uname] = upass
	mu.Unlock()
}

func isValidUser(uname string, upass string) bool {
	mu.RLock()
	val, isValid := userData[uname]
	mu.RUnlock()
	if !isValid {
		return false
	}
	return val == upass
}

func registerPage(c *gin.Context){
	c.HTML(200, "register.html", nil)
}

func registerAPI(c *gin.Context) {
	uname := c.PostForm("username")
	upass := c.PostForm("password")
	if isValidUser(uname, upass) {
		c.String(http.StatusConflict, "Username already exists!")
		time.Sleep(2*time.Second)
		c.Redirect(http.StatusMovedPermanently, "register")
		return
	}
	mu.RLock()
	userData[uname] = upass
	usernameToLinks[uname] = make([]string, 0, 0)
	c.SetCookie("username", uname, 3600, "/", "", false, true)
	c.SetCookie("password", upass, 3600, "/", "", false, true)
	mu.RUnlock()
	c.Redirect(http.StatusMovedPermanently, "user")
}

func loginPage(c *gin.Context){
	c.HTML(200, "login.html", nil)
}

func loginAPI(c *gin.Context){
	uname := c.PostForm("username")
	upass := c.PostForm("password")
	if !isValidUser(uname, upass){
		c.String(400, "Invalid credentials!")
		time.Sleep(2*time.Second)
		c.Redirect(http.StatusMovedPermanently, "login")
	}
}

func mainPage(c *gin.Context){
	uname, err1 := c.Cookie("username")
	upass, err2 := c.Cookie("password")
	if err1 != nil || err2 != nil {
		c.HTML(200, "login.html", nil)
		return
	}
	if !isValidUser(uname, upass){
		c.HTML(http.StatusUnauthorized, "login.html", nil)
		return
	}
	mu.RLock()
	sli, _ := usernameToLinks[uname]
	c.HTML(200, "user.html", gin.H{
		"username" : uname,
		"Items" : sli,
	})
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("static/*")
	router.GET("/", mainPage)
	router.GET("/register", registerPage)
	router.POST("/register", registerAPI)
	router.GET("/login", loginPage)
	router.POST("/login", loginAPI)
	router.Run()
}

/*obsolete
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
			if userData[uname] == pass {
				c.SetCookie("username", uname, 3600, "/", "", true, true)
				c.SetCookie("password", pass, 3600, "/", "", true, true)
				c.Redirect(http.StatusPermanentRedirect, "/user/"+uname)
			} else {
				c.String(http.StatusConflict, "Bad Login.")
			}
			mu.RUnlock()
		} else {
			c.String(http.StatusUnauthorized, "Bad Credentials.")
		}
	})

	r.GET("/user/:username", func(c *gin.Context) {
		uname := c.Param("username")
		cache_name, err := c.Cookie("username")
		if err != nil {
			c.String(http.StatusUnauthorized, "Can't view without login")
			return
		}
		if uname != cache_name {
			c.String(http.StatusUnauthorized, "Can't view without login")
			return
		}
		mu.RLock()
		val, isValid := userData[uname]
		cache_pass, err1 := c.Cookie("password")
		if err1 != nil || !isValid || val != cache_pass {
			mu.RLock()
			c.String(http.StatusUnauthorized, "Can't view without login")
			return
		}
		slice, exists := usernameToLinks[uname]
		res := "Login as "+uname + "\n\n" + "Files:\n"
		if exists{
			for _, x := range(slice) {
				res = res + "\t" + x + "\n"
			}
		}
		c.String(200, res)
		
		mu.RUnlock()
	})

	r.POST("/upload", func(c *gin.Context) {
		uname, err1 := c.Cookie("username")
		pass, err2 := c.Cookie("password")
		if err1 != nil || err2 != nil || !isValidUser(uname, pass) {
			c.String(http.StatusUnauthorized, "Bad Credentials")
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "File could not be opened.")
			return
		}
		c.SaveUploadedFile(file, "uploads/"+uname+"/"+file.Filename)
		path := "uploads/"+uname+"/"+file.Filename
		c.String(http.StatusOK, "File Successfuly Uploaded to: "+path)
		_, exists := usernameToLinks[uname]
		if !exists{
			usernameToLinks[uname] = make([]string, 1, 1)
			usernameToLinks[uname][0] = path
		}else{
			usernameToLinks[uname] = append(usernameToLinks[uname], path)
		}
	})
	r.GET("/uploads/:username/:filename", func(c *gin.Context) {
		uname, err1 := c.Cookie("username")
		pass, err2 := c.Cookie("password")
		if err1 != nil || err2 != nil {
			c.String(http.StatusUnauthorized, "Please Login.")
			return
		}
		paramUserName := c.Param("username")
		if uname != paramUserName {
			c.String(http.StatusUnauthorized, "Bad Login.")
			return
		}
		if !isValidUser(uname, pass) {
			c.String(http.StatusUnauthorized, "Bad Login.")
			return
		}
		paramFileName := c.Param("filename")
		c.File("./uploads/" + uname + "/" + paramFileName)
	})
	r.GET("/fetchUploadList", func(c *gin.Context){
		uname, err1 := c.Cookie("username")
		upass, err2 := c.Cookie("password")
		if err1 != nil || err2 != nil {
			c.String(http.StatusUnauthorized, "Please Login.")
			return
		}
		if !isValidUser(uname, upass) {
			c.String(http.StatusUnauthorized, "Bad Login.")
			return
		}
		mu.RLock()
		sli, valid := usernameToLinks[uname]
		mu.RUnlock()
		if !valid{
			c.JSON(http.StatusOK, nil)
			return
		}
		c.JSON(http.StatusOK, sli)
		
	})
	r.Run()
}
*/