package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arismonx/Coffee-Tuna-Liff/config"
	"github.com/google/generative-ai-go/genai"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig(".env")

	// ==========================================
	// 1. ตั้งค่าการเชื่อมต่อ (Gemini & Weaviate)
	// ==========================================
	aiClient, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		log.Fatalf("เชื่อมต่อ Gemini ไม่สำเร็จ: %v", err)
	}
	defer aiClient.Close()

	wCfg := weaviate.Config{Host: "localhost:8080", Scheme: "http"}
	wClient, err := weaviate.NewClient(wCfg)
	if err != nil {
		log.Fatalf("เชื่อมต่อ Weaviate ไม่สำเร็จ: %v", err)
	}

	// ==========================================
	// 2. รับคำถามและแปลงเป็น Vector
	// ==========================================
	question := "ร้านเปิดกี่โมงครับ?"
	fmt.Printf("👤 คำถามผู้ใช้: %s\n", question)

	em := aiClient.EmbeddingModel("models/gemini-embedding-001")
	embedRes, err := em.EmbedContent(ctx, genai.Text(question))
	if err != nil {
		log.Fatalf("แปลง Vector คำถามไม่สำเร็จ: %v", err)
	}
	questionVector := embedRes.Embedding.Values

	// ==========================================
	// 3. ค้นหาข้อมูลที่ใกล้เคียงที่สุดใน Weaviate
	// ==========================================
	// กำหนดเงื่อนไขการค้นหา (ค้นหาด้วย Vector และดึงมาแค่ 1 รายการที่ใกล้เคียงสุด)
	nearVec := wClient.GraphQL().NearVectorArgBuilder().WithVector(questionVector)

	result, err := wClient.GraphQL().Get().
		WithClassName("Document").
		WithFields(graphql.Field{Name: "text"}).
		WithNearVector(nearVec).
		WithLimit(1).
		Do(ctx)

	if err != nil {
		log.Fatalf("ค้นหาใน Weaviate ไม่สำเร็จ: %v", err)
	}

	// แกะเอาข้อความที่ค้นเจอออกมาจาก JSON ของ Weaviate
	var retrievedText string
	data := result.Data
	if get, ok := data["Get"].(map[string]interface{}); ok {
		if docs, ok := get["Document"].([]interface{}); ok && len(docs) > 0 {
			if doc, ok := docs[0].(map[string]interface{}); ok {
				retrievedText = doc["text"].(string)
			}
		}
	}
	fmt.Printf("\n🔍 ข้อมูลที่บอทค้นเจอจาก Database: %s\n", retrievedText)

	// ==========================================
	// 4. ส่งข้อมูลให้ Gemini สรุปเป็นคำตอบ (Prompt Engineering)
	// ==========================================
	// เราจะเอาข้อมูลที่ค้นเจอ มาบวกกับคำถาม แล้วสั่งให้ AI ตอบจากข้อมูลนี้เท่านั้น
	prompt := fmt.Sprintf(`คุณคือพนักงานร้านกาแฟ ตอบคำถามลูกค้าโดยใช้ "ข้อมูลอ้างอิง" ที่กำหนดให้เท่านั้น ห้ามเดาเอาเอง
    
ข้อมูลอ้างอิง: %s

คำถามลูกค้า: %s`, retrievedText, question)

	chatModel := aiClient.GenerativeModel("gemini-2.5-flash")
	chatRes, err := chatModel.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatalf("Gemini สร้างคำตอบไม่สำเร็จ: %v", err)
	}

	// แสดงผลคำตอบ
	for _, part := range chatRes.Candidates[0].Content.Parts {
		fmt.Printf("\n🤖 บอทตอบกลับ:\n%v\n", part)
	}
}
