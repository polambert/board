
package main

import (
	_ "fmt"

	"net/http"

	"github.com/gin-gonic/gin"

	"time"

	"strconv"

	"io/ioutil"

	"encoding/json"
)

type Comment struct {
	Time time.Time
	TimeString string
	Body string
	Country string
	Id int
}

type Post struct {
	Time time.Time
	TimeString string
	Body string
	Country string
	Id int
	Comments []Comment
}

func main() {
	posts := make([]Post, 0)
	var id int = 0

	router := gin.Default()

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"posts": posts,
		})
	})

	router.GET("/post/:id", func(c *gin.Context) {
		post := interface{}(nil)

		for i := 0; i < len(posts); i++ {
			id, _ := strconv.Atoi(c.Param("id"))
			if posts[i].Id == id {
				post = posts[i]
				break
			}
		}

		c.HTML(http.StatusOK, "post.html", gin.H{
			"post": post,
		})
	})

	router.POST("/post", func(c *gin.Context) {
		// Create the post.
		t := time.Now()
		ts := t.Format("Mon Jan 2 15:04:05 PM")

		resp, _ := http.Get("http://api.ipstack.com/" + c.Request.RemoteAddr + "?access_key=2afa2c93b8578e6d891697f56cc8ccdf")
		defer resp.Body.Close()
		text, _ := ioutil.ReadAll(resp.Body)

		jsondata := make(map[string]interface{})
		json.Unmarshal(text, &jsondata)

		country := jsondata["country_name"]
		country_name := "Unknown Country"

		if country != nil {
			country_name = country.(string)
		}

		post := Post{
			t,
			ts,
			c.PostForm("body"),
			country_name,
			id,
			make([]Comment, 0),
		}

		posts = append(posts, post)

		// Redirect the user to the post page.
		c.Redirect(http.StatusMovedPermanently, "/post/" + strconv.Itoa(id))
		c.Abort()

		id += 1
	})

	router.POST("/comment/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		t := time.Now()
		ts := t.Format("Mon Jan 2 15:04:05 PM")

		resp, _ := http.Get("http://api.ipstack.com/" + c.Request.RemoteAddr + "?access_key=2afa2c93b8578e6d891697f56cc8ccdf")
		defer resp.Body.Close()
		text, _ := ioutil.ReadAll(resp.Body)

		jsondata := make(map[string]interface{})
		json.Unmarshal(text, &jsondata)

		country := jsondata["country_name"]
		country_name := "Unknown Country"

		if country != nil {
			country_name = country.(string)
		}

		for i := 0; i < len(posts); i++ {
			if posts[i].Id == id {
				posts[i].Comments = append(posts[i].Comments, Comment{
					t,
					ts,
					c.PostForm("body"),
					country_name,
					len(posts[i].Comments),
				})

				break
			}
		}

		c.Redirect(http.StatusMovedPermanently, "/post/" + c.Param("id"))
		c.Abort()
	})

	router.Run(":80")
}
