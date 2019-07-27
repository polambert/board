
package main

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting bitter.polambert.com.")

	router := gin.Default()

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	router.Run(":80")
}
