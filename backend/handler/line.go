package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Arismonx/Coffee-Tuna-Liff/config"
	"github.com/Arismonx/Coffee-Tuna-Liff/model"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
)

// this past is format for Events Webhook from request line server
// more detail https://developers.line.biz/en/reference/messaging-api/#webhooks

// strcut events
type EventsPayload struct {
	Events []WebhookEvent
}

type WebhookEvent struct {
	Type       string        `json:"type"`
	ReplyToken string        `json:"replyToken"`
	Message    EventMessage  `json:"message"`
	Postback   EventPostback `json:"postback"`
}

type EventPostback struct {
	Data string `json:"data"`
}

type EventMessage struct {
	Type            string `json:"type"`
	Text            string `json:"text"`
	MarkAsReadToken string `json:"markAsReadToken"`
}

// struct JSON Reply Payload send to line server
type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ReplyPayload struct {
	ReplyToken string    `json:"replyToken"`
	Messages   []Message `json:"messages"`
}

func CreateSendRequesReply(payload ReplyPayload, Token string, url string) {
	// convert struct to string JSON (look like JSON.stringify)
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	// Create Requset to Line Server
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	// Set Header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Token)

	// Sent Request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}

	// Close
	defer res.Body.Close()

	// Check LINE Response
	log.Println("LINE Response Status:", res.Status)
}

// strcut config
type LineHandler struct {
	Config config.Config
	Model  *genai.GenerativeModel
}

// Construct LineHandler
func NewLineHandler(
	Config config.Config,
	Model *genai.GenerativeModel,
) *LineHandler {
	return &LineHandler{
		Config: Config,
		Model:  Model,
	}
}

// Bind this func with struct LineHandler and can use attribute in struct
func (h *LineHandler) Webhook(ctx *gin.Context) {

	// Line Channel Access Token
	token := h.Config.LineChannelAccessToken

	// URL Reply
	url := "https://api.line.me/v2/bot/message/reply"

	// object event data webhook from request line server
	var body EventsPayload

	if err := ctx.ShouldBindJSON(&body); err != nil {
		fmt.Println("Error Bind JSON:", err)
		ctx.JSON(400, gin.H{"error": err})
		return
	}

	// Check events data have object events / Chat with Ai
	for _, event_0 := range body.Events {

		user_message := event_0.Message.Text
		fmt.Println("User Message: ", user_message)

		// Check type is type message
		if event_0.Type == "message" {

			readToken := event_0.Message.MarkAsReadToken

			// === Read status ===
			markReadPayload := map[string]string{"markAsReadToken": readToken}
			jsonData, _ := json.Marshal(markReadPayload)
			req, _ := http.NewRequest("POST", "https://api.line.me/v2/bot/chat/markAsRead", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			client := &http.Client{}
			client.Do(req)
			// =========================================================

			ctx_ai := context.Background()
			ctx_ai, cancel := context.WithTimeout(ctx_ai, 20*time.Second)
			defer cancel()

			// Call function GenerateContent_textOnly get Text
			resp_message := model.GenerateContent_textOnly(ctx_ai, user_message, h.Model)

			// Payload Data Reply
			payload := ReplyPayload{
				ReplyToken: event_0.ReplyToken,
				Messages: []Message{
					{Type: "text", Text: resp_message},
				},
			}

			CreateSendRequesReply(payload, token, url)

		}

		// Check type is type postback
		if event_0.Type == "postback" {
			fmt.Println("Data: ", event_0.Postback.Data)
			if event_0.Postback.Data == "ปุ่มAนะ" {
				payload := ReplyPayload{
					ReplyToken: event_0.ReplyToken,
					Messages: []Message{
						{Type: "text", Text: "ปุ่มAนะ"},
					},
				}

				CreateSendRequesReply(payload, token, url)
			}
		}

	}
	ctx.JSON(200, gin.H{"message": "HTTP POST request sent to the webhook URL!"})
}
