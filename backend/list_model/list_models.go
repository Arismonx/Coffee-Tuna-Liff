package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Arismonx/Coffee-Tuna-Liff/config"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Load .ENV
	cfg := config.LoadConfig(".env")

	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		log.Fatalf("สร้าง Client ไม่สำเร็จ: %v", err)
	}
	defer client.Close()

	fmt.Println("ค้นหารายชื่อ Embedding Models ที่คุณใช้งานได้...")

	iter := client.ListModels(ctx)
	found := false

	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("เกิดข้อผิดพลาดตอนดึงข้อมูลโมเดล: %v", err)
		}

		// กรองดูเฉพาะโมเดลที่มีคำว่า "embed" ในชื่อ
		if strings.Contains(m.Name, "embed") {
			fmt.Printf("=> ค้นพบโมเดล: %s\n", m.Name)
			found = true
		}
	}

	if !found {
		fmt.Println("ไม่พบโมเดลสำหรับทำ Embedding ใน API Key นี้เลย")
	}
}
