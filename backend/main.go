package main

import (
	"fmt"
	"net/http"

	"github.com/Arismonx/Coffee-Tuna-Liff/config"
	"github.com/Arismonx/Coffee-Tuna-Liff/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	// load ENV path and specified file .env
	cfg := config.LoadConfig(".env")

	// create LineHandler and assing cfg to attribute Config
	lineHandler := handler.NewLineHandler(cfg)

	router := gin.Default()

	// == API ==
	// router "/" is mean assume Home Page
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Home Page"})
	})

	// router "/webhook" is Send Message into Line OA / Message API
	router.POST("/webhook", lineHandler.Webhook)

	// Run Prot 8000
	fmt.Printf("Start : http://localhost:8000 \n")
	router.Run(":8000")
}
