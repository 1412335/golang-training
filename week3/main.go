package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	counter := 0
	r.GET("/increase", func(c *gin.Context) {
		counter++
		c.JSON(200, gin.H{
			"counter": counter,
		})
	})
	r.Run()
}
