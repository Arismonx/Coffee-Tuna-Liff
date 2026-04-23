package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Arismonx/Coffee-Tuna-Liff/config"
	"github.com/Arismonx/Coffee-Tuna-Liff/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {

	ctx := context.Background()

	// load ENV path and specified file .env
	cfg := config.LoadConfig(".env")

	// Create Model Ai Gemini
	// model_gemini := model.CreateModel(ctx, cfg.GeminiAPIKey, "gemini-2.5-flash")

	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		log.Printf("Gemini Error NewClient: %v", err)
	}
	defer client.Close()

	// send to model "gemini-2.5-flash"
	model := client.GenerativeModel("gemini-2.5-flash")
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text("คุณคือบอทตอบคำถามของร้านกาแฟ ให้ตอบแบบเป็นกันเอง สั้น กระชับ และตรงประเด็นที่สุด ไม่ต้องเกริ่นนำเยอะ"),
		},
	}

	// create LineHandler and assing cfg to attribute Config
	lineHandler := handler.NewLineHandler(cfg, model)

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
