package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	fmt.Println("Running App")
	//Tomar archivos
	router.LoadHTMLGlob("templates/*")
	router.GET("/ping", func(c *gin.Context) { //se define una url
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	router.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
