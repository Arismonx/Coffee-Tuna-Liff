package main

import (
	"context"
	"fmt"
	"log"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	ctx := context.Background()

	// 1. ตั้งค่าเชื่อมต่อ Weaviate บนเครื่อง Local
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Fatalf("เชื่อมต่อ Weaviate ไม่สำเร็จ: %v", err)
	}

	// 2. ลบ Class เก่าทิ้ง (ถ้ามี) เพื่อให้รันซ้ำได้ง่าย
	className := "Document"
	client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

	// 3. สร้าง Class (เหมือนการสร้าง Table ใน SQL)
	classObj := &models.Class{
		Class: className,
		Properties: []*models.Property{
			{
				Name:     "text",
				DataType: []string{"text"},
			},
		},
	}
	err = client.Schema().ClassCreator().WithClass(classObj).Do(ctx)
	if err != nil {
		log.Fatalf("สร้าง Class ไม่สำเร็จ: %v", err)
	}
	fmt.Println("=> สร้าง Class 'Document' สำเร็จ")

	// 4. เตรียมข้อมูลจำลอง (ในระบบจริง Vector นี้จะได้มาจาก Gemini Embedding)
	mockText := "ร้าน Coffee Tuna เปิดให้บริการทุกวัน เวลา 08:00 - 18:00 น."
	// สมมติว่านี่คือ Vector ที่ได้จาก Gemini (ปกติต้องมี 768 ตัว)
	mockVector := []float32{0.1, 0.2, -0.3, 0.4, 0.5}

	// 5. บันทึกข้อมูลลง Weaviate
	wrapper, err := client.Data().Creator().
		WithClassName(className).
		WithProperties(map[string]interface{}{
			"text": mockText,
		}).
		WithVector(mockVector). // แนบ Vector เข้าไปด้วย
		Do(ctx)

	if err != nil {
		log.Fatalf("บันทึกข้อมูลไม่สำเร็จ: %v", err)
	}

	fmt.Printf("=> บันทึกข้อมูลสำเร็จ! ID ของข้อมูลคือ: %s\n", wrapper.Object.ID)
}
