package main

import "github.com/gin-gonic/gin"

type name interface {
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run("0.0.0.0:3000") // listens on 0.0.0.0:8080 by default
}
