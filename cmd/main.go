package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
