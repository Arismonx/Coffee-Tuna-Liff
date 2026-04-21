package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// strcut events
type EventsPayload struct {
	Events []WebhookEvent
}

type WebhookEvent struct {
	Type       string       `json:"type"`
	ReplyToken string       `json:"replyToken"`
	Message    EventMessage `json:"message"`
}

type EventMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// struct JSON Payload send to line server
type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ReplyPayload struct {
	ReplyToken string    `json:"replyToken"`
	Messages   []Message `json:"messages"`
}

func webhook(ctx *gin.Context) {
	Token := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	// URL Reply
	url := "https://api.line.me/v2/bot/message/reply"

	// object even data webhook
	var body EventsPayload

	if err := ctx.ShouldBindJSON(&body); err != nil {
		fmt.Println("Error Bind JSON:", err)
		ctx.JSON(400, gin.H{"error": "รูปแบบ JSON ไม่ถูกต้อง"})
		return
	}

	// Check events data have object events
	if len(body.Events) > 0 {
		event_0 := body.Events[0]

		if event_0.Type == "message" {
			// Payload Data Reply
			payload := ReplyPayload{
				ReplyToken: event_0.ReplyToken,
				Messages: []Message{
					{Type: "text", Text: "Hello World"},
				},
			}

			// convert struct to string JSON (look like JSON.stringify)
			jsonData, err := json.Marshal(payload)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				return
			}

			// Create Requset to Line Server
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			// Set Header
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+Token)

			// Sent Request
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}

			// Close
			defer req.Body.Close()

			// Check LINE Response
			fmt.Println("LINE Response Status:", res.Status)
		}

	}

	ctx.JSON(200, gin.H{"message": "HTTP POST request sent to the webhook URL!"})
}

func main() {
	// load ENV
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	// == API ==
	// router "/" is mean assume Home Page
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Home Page"})
	})

	// router "/webhook" is Send Message into Line OA / Message API
	router.POST("/webhook", webhook)

	// Run Prot 8000
	fmt.Printf("Start : http://localhost:8000 \n")
	router.Run(":8000")
}
