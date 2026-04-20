package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// Api
	// router "/" is mean assume Home Page
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Home Page"})
	})

	// router "/send-message" is Send Message into Line OA / Message API
	router.POST("/send-message", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Run Prot 8000
	fmt.Printf("Start : http://localhost:8000 \n")
	router.Run(":8000")
}
