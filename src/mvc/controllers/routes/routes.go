package routes

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

type Person struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func SetupRouter(router *gin.Engine) *gin.Engine {

	// -----------------------------------
	//  Public
	// -----------------------------------
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "public/index.html", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.html", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.html", gin.H{
			"title": "Users",
		})
	})

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	router.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// For example :http://localhost:5000/person/somname/c974912f-5a4d-42a8-862a-893644f3bb6a
	router.GET("/person/:name/:id", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			log.Printf("Error --> %v", err)
			errMsg := fmt.Sprintf("%s", err)
			c.JSON(400, gin.H{"msg": errMsg})
			return
		}
		c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
	})

	// Goroutine @docs https://gin-gonic.com/docs/examples/goroutines-inside-a-middleware/
	router.GET("/long_async", func(c *gin.Context) {
		// create copy to be used inside the goroutine
		cCp := c.Copy()
		go func() {
			// simulate a long task
			time.Sleep(2 * time.Second)

			// note that you are using the copied context "cCp", IMPORTANT
			log.Println("Done! in path " + cCp.Request.URL.Path)
			for i := 0; i < 100; i++ {
				log.Printf("GoRoutine Printing ---> %d\n", i)
			}
		}()

		c.JSON(200, gin.H{"page": "long_async"})
	})

	// -----------------------------------
	// 	V1
	// -----------------------------------
	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/login", func(c *gin.Context) { c.JSON(200, gin.H{"page": "V1 login"}) })
		v1.GET("/submit", func(c *gin.Context) { c.JSON(200, gin.H{"page": "V1 submit"}) })
		v1.GET("/read", func(c *gin.Context) { c.JSON(200, gin.H{"page": "V1 read"}) })
	}

	// -----------------------------------
	// 	V2
	// -----------------------------------
	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.GET("/login", func(c *gin.Context) { c.JSON(200, gin.H{"page": "V2 login"}) })
		v2.GET("/submit", func(c *gin.Context) { c.JSON(200, gin.H{"page": "V2 submit"}) })
		v2.GET("/read", func(c *gin.Context) { c.JSON(200, gin.H{"page": "V2 read"}) })
	}

	// -----------------------------------
	// 	Logged In
	// -----------------------------------
	authorizedAccounts := gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}

	logged := router.Group("/", gin.BasicAuth(authorizedAccounts))

	logged.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		fmt.Println(user)
		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Aunothorized"})
	})

	return router

}
