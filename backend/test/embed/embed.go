package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Arismonx/Coffee-Tuna-Liff/config"
	"github.com/google/generative-ai-go/genai"
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

	// 2. เรียกใช้โมเดล Embedding
	em := client.EmbeddingModel("models/gemini-embedding-001")

	// 3. เปิดไฟล์ data.txt
	file, err := os.Open("docs/menu.txt")
	if err != nil {
		log.Fatalf("เปิดไฟล์ไม่สำเร็จ: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 1

	// 4. อ่านทีละบรรทัด
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue // ข้ามบรรทัดว่าง
		}

		fmt.Printf("บรรทัดที่ %d: %s\n", lineCount, text)

		// 5. ส่งข้อความไปแปลงเป็น Vector (ชุดตัวเลข)
		res, err := em.EmbedContent(ctx, genai.Text(text))
		if err != nil {
			log.Printf("เกิดข้อผิดพลาดในการแปลง Vector: %v", err)
			continue
		}

		// ตัวแปร vector นี้คือ []float32 ที่เก็บตัวเลขไว้
		vector := res.Embedding.Values

		// แสดงผลขนาดของ Vector และตัวอย่างตัวเลข 3 ตัวแรก
		fmt.Printf("=> สำเร็จ! ได้ Vector ขนาด %d มิติ\n", len(vector))
		fmt.Printf("=> ตัวอย่างค่า: [%f, %f, %f, ...]\n\n", vector[0], vector[1], vector[2])

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("เกิดข้อผิดพลาดในการอ่านไฟล์: %v", err)
	}
}
