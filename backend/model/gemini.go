package model

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func CreateModel(ctx context.Context, token string, seleteModel string) *genai.GenerativeModel {
	// create request
	client, err := genai.NewClient(ctx, option.WithAPIKey(token))
	if err != nil {
		log.Printf("Gemini Error NewClient: %v", err)
	}
	defer client.Close()

	// send to model "gemini-2.5-flash"
	model := client.GenerativeModel(seleteModel)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text("คุณคือบอทตอบคำถามของร้านกาแฟ ให้ตอบแบบเป็นกันเอง สั้น กระชับ และตรงประเด็นที่สุด ไม่ต้องเกริ่นนำเยอะ"),
		},
	}
	return model
}

// ============================================================================
// Copyright 2023 Google LLC
// Copyright 2026 Arismonx

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ============================================================================

func GenerateContent_textOnly(ctx context.Context, message string, model *genai.GenerativeModel) string {
	// response
	resp, err := model.GenerateContent(ctx, genai.Text(message))
	if err != nil {
		log.Printf("Gemini Error: %v", err)
		return "Error: System crash at response"
	}

	return Response_TextFromGenerative_Model(resp)
}

func Response_TextFromGenerative_Model(resp *genai.GenerateContentResponse) string {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if text, ok := part.(genai.Text); ok {
					log.Println(text)
					return string(text)
				}
			}
		}
	}
	fmt.Println("---")
	return "Error: empty string"
}
