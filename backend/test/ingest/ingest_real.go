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
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig(".env")

	// 1. เชื่อมต่อ Gemini และ Weaviate
	aiClient, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		log.Fatalf("เชื่อมต่อ Gemini ไม่สำเร็จ: %v", err)
	}
	defer aiClient.Close()

	wClient, err := weaviate.NewClient(weaviate.Config{Host: "localhost:8080", Scheme: "http"})
	if err != nil {
		log.Fatalf("เชื่อมต่อ Weaviate ไม่สำเร็จ: %v", err)
	}

	className := "Document"

	// 2. ล้าง Class เก่าทิ้งแล้วสร้างใหม่
	wClient.Schema().ClassDeleter().WithClassName(className).Do(ctx)
	classObj := &models.Class{
		Class: className,
		Properties: []*models.Property{
			{Name: "text", DataType: []string{"text"}},
		},
	}
	err = wClient.Schema().ClassCreator().WithClass(classObj).Do(ctx)
	if err != nil {
		log.Fatalf("สร้าง Class ไม่สำเร็จ: %v", err)
	}

	// 3. เปิดไฟล์ docs/menu.txt
	file, err := os.Open("docs/menu.txt")
	if err != nil {
		log.Fatalf("เปิดไฟล์ docs/menu.txt ไม่สำเร็จ: %v", err)
	}
	defer file.Close()

	em := aiClient.EmbeddingModel("models/gemini-embedding-001")
	scanner := bufio.NewScanner(file)

	// 4. อ่านทีละบรรทัด แปลงเป็น Vector และบันทึกลง Database
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		// แปลงข้อความบรรทัดนี้เป็น Vector 768 มิติ
		res, err := em.EmbedContent(ctx, genai.Text(text))
		if err != nil {
			log.Printf("แปลง Vector ไม่ผ่าน: %v", err)
			continue
		}
		vector := res.Embedding.Values

		// บันทึกลง Weaviate
		_, err = wClient.Data().Creator().
			WithClassName(className).
			WithProperties(map[string]interface{}{
				"text": text,
			}).
			WithVector(vector).
			Do(ctx)

		if err != nil {
			log.Printf("บันทึกข้อมูลไม่สำเร็จ: %v", err)
		} else {
			fmt.Printf("✅ บันทึกสำเร็จ: %s\n", text)
		}
	}
	fmt.Println("\n🎉 อัปโหลดข้อมูลจริงเข้า Database ทั้งหมดเรียบร้อยแล้ว!")
}
