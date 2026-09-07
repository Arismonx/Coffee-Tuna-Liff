# Coffee Tuna LIFF — LINE OA Chatbot (RAG)

แชทบอทสำหรับ **LINE Official Account** ของร้านกาแฟ "Coffee Tuna"
ตอบคำถามลูกค้าแบบ **RAG (Retrieval-Augmented Generation)** โดยใช้

- **Go + Gin** เป็น Webhook server รับ event จาก LINE Messaging API
- **Google Gemini** (`gemini-2.5-flash`) สร้างคำตอบ + (`gemini-embedding-001`) แปลงข้อความเป็น vector
- **Weaviate** เป็น vector database เก็บข้อมูลร้าน (เมนู, เวลาเปิด-ปิด, โปรโมชั่น ฯลฯ)
- **LIFF** หน้าเว็บ static แสดงโปรไฟล์ผู้ใช้ที่ login ผ่าน LINE

---

## สถาปัตยกรรม

```
                ┌──────────────┐   Webhook (POST /webhook)   ┌────────────────────┐
  ผู้ใช้ LINE ──▶│  LINE Server │ ──────────────────────────▶ │  Go Backend (Gin)  │
                └──────────────┘                              │      :8000          │
                       ▲                                      └─────────┬──────────┘
                       │  Reply API                                     │
                       │                             ┌──────────────────┼───────────────────┐
                       │                             ▼                  ▼                   ▼
                       │                    embed คำถาม (Gemini)  ค้น vector (Weaviate)  สร้างคำตอบ (Gemini)
                       └───────────────────────────── ข้อความตอบกลับ ◀──────────────────────┘

  LIFF page (static, :8888) ──▶ liff.init / liff.getProfile ──▶ LINE Login
```

**ขั้นตอนเมื่อมีข้อความเข้ามา** (`backend/handler/line.go`):

1. รับ event จาก `POST /webhook`
2. เรียก `chat/markAsRead` ทำเครื่องหมายว่าอ่านแล้ว
3. แปลงข้อความผู้ใช้เป็น vector 768 มิติด้วย `gemini-embedding-001`
4. ค้นข้อมูลใกล้เคียงที่สุด 1 รายการจาก Weaviate class `Document` (near vector)
5. ประกอบ RAG prompt = ข้อมูลอ้างอิง + คำถามลูกค้า แล้วส่งให้ `gemini-2.5-flash`
6. ตอบกลับผ่าน LINE Reply API (`/v2/bot/message/reply`)
7. ถ้าเป็น event ชนิด `postback` และ `data == "ปุ่มAนะ"` จะตอบข้อความคงที่กลับไป

---

## โครงสร้างโปรเจกต์

```
lineOA/
├── backend/
│   ├── main.go                     # entry point, ตั้งค่า Gin router + client ต่าง ๆ, port :8000
│   ├── docker-compose.yml          # Weaviate service (:8080 REST, :50051 gRPC)
│   ├── .env.example                # ตัวอย่างตัวแปรสภาพแวดล้อม
│   ├── config/
│   │   ├── config.go               # โหลด .env -> struct Config
│   │   └── config_test.go
│   ├── handler/
│   │   └── line.go                 # LineHandler.Webhook — logic หลักของบอท
│   ├── model/
│   │   └── gemini.go               # สร้าง GenerativeModel + generate ข้อความ
│   ├── list_model/
│   │   └── list_models.go          # utility: list embedding models ที่ API key ใช้ได้
│   ├── docs/
│   │   └── menu.txt                # knowledge base (1 บรรทัด = 1 เอกสาร)
│   └── test/                       # สคริปต์ทดสอบ/เครื่องมือ (แต่ละไฟล์เป็น package main แยกกัน)
│       ├── ingest/ingest_real.go   # อ่าน menu.txt -> embed -> เขียนลง Weaviate
│       ├── embed/embed.go          # ทดสอบการ embed ทีละบรรทัด
│       └── rag/main.go             # ทดสอบ RAG pipeline แบบ end-to-end (คำถาม hardcode)
└── frontend/
    ├── index.html                  # หน้า LIFF (LIFF SDK 2.22.3)
    └── package.json                # script: serve static ผ่าน python http.server :8888
```

---

## สิ่งที่ต้องมีก่อนเริ่ม

| เครื่องมือ | เวอร์ชันแนะนำ | ใช้ทำอะไร |
|---|---|---|
| Go | 1.25+ (ตาม `backend/go.mod`) | build/run backend |
| Docker + Docker Compose | ล่าสุด | รัน Weaviate |
| Python 3 | 3.x | serve หน้า LIFF (static) |
| Gemini API Key | — | เอาจาก <https://aistudio.google.com/apikey> |
| LINE Official Account + Messaging API channel | — | เอา Channel Access Token, ตั้ง Webhook |
| ngrok หรือ public HTTPS domain | — | เปิด `/webhook` ให้ LINE เรียกเข้ามาได้ตอน dev |

---

## การติดตั้ง

### 1. Clone และตั้งค่า environment

```bash
git clone <repo-url>
cd lineOA/backend
cp .env.example .env
```

แก้ไฟล์ `backend/.env`:

```env
LINE_CHANNEL_ACCESS_TOKEN="<channel access token จาก LINE Developers Console>"
GEMINI_API_KEY=<gemini api key>
```

> `config.LoadConfig(".env")` โหลดไฟล์แบบ path สัมพัทธ์ ดังนั้นทุกคำสั่ง `go run` ต้องรันจากไดเรกทอรี `backend/`

### 2. ติดตั้ง dependencies ของ backend

```bash
cd backend
go mod download
```

### 3. รัน Weaviate (vector database)

```bash
cd backend
docker compose up -d
```

- REST API: `http://localhost:8080`
- ข้อมูลถูก persist ไว้ที่ `backend/weaviate_data/`
- ปิด/หยุด: `docker compose down`

### 4. ป้อนข้อมูลร้านเข้า Weaviate (ingest)

แก้/เพิ่มข้อมูลใน `backend/docs/menu.txt` (1 บรรทัด = 1 เอกสาร) แล้วรัน:

```bash
cd backend
go run ./test/ingest/
```

สคริปต์นี้จะ **ลบ class `Document` เดิมทิ้งแล้วสร้างใหม่** จากนั้น embed ทุกบรรทัดและบันทึกลง Weaviate
รันซ้ำได้ทุกครั้งที่แก้ `menu.txt`

---

## การรันโปรแกรม

### รัน Backend (Webhook server)

```bash
cd backend
go run main.go
```

- เปิดที่ `http://localhost:8000`
- `GET  /`         → health check (`{"message":"Home Page"}`)
- `POST /webhook`  → endpoint สำหรับ LINE Messaging API

build เป็น binary:

```bash
cd backend
go build -o coffee-tuna .
./coffee-tuna
```

### รัน Frontend (หน้า LIFF)

```bash
cd frontend
npm run client
```

เปิดที่ `http://localhost:8888` (เสิร์ฟไฟล์ static ผ่าน `python3 -m http.server 8888`)
หน้านี้จะเรียก `liff.init()` ด้วย LIFF ID ที่ hardcode อยู่ใน `frontend/index.html`
ถ้า login แล้วจะแสดงรูปโปรไฟล์ ชื่อ และ userId

---

## เครื่องมือ / สคริปต์ทดสอบ

รันทั้งหมดจากไดเรกทอรี `backend/`

| คำสั่ง | หน้าที่ |
|---|---|
| `go run ./test/ingest/` | อ่าน `docs/menu.txt` → embed → เขียนลง Weaviate (สร้าง class ใหม่) |
| `go run ./test/embed/` | ทดสอบการแปลงข้อความเป็น vector ทีละบรรทัดจาก `docs/menu.txt` |
| `go run ./test/rag/` | ทดสอบ RAG pipeline ครบวงจร (คำถามตัวอย่าง `"ร้านเปิดกี่โมงครับ?"`) — ต้อง ingest ข้อมูลก่อน |
| `go run ./list_model/` | แสดงรายชื่อ embedding model ที่ Gemini API key ใช้งานได้ |
| `go test ./config/ -v` | unit test ของ `LoadConfig` (ต้องมี `backend/.env` ที่มีค่าจริง) |

---

## การเชื่อมต่อ LINE Webhook

### 1. สร้าง Messaging API channel

1. ไปที่ [LINE Developers Console](https://developers.line.biz/console/) → สร้าง **Provider** → สร้าง channel ชนิด **Messaging API**
2. ที่แท็บ **Messaging API**
   - กด **Issue** ที่ *Channel access token (long-lived)* → นำค่าไปใส่ `LINE_CHANNEL_ACCESS_TOKEN` ใน `.env`
   - ปิด **Auto-reply messages** และ **Greeting messages** (ไม่งั้นจะตอบชนกับบอท)

### 2. เปิด backend ออกสู่อินเทอร์เน็ต (ตอน dev)

LINE ต้องเรียก webhook ผ่าน **HTTPS สาธารณะ** ใช้ ngrok ครอบ port 8000:

```bash
ngrok http 8000
```

จะได้ URL เช่น `https://xxxx-xxx-xxx.ngrok-free.app`

### 3. ตั้งค่า Webhook URL

ที่แท็บ **Messaging API** ของ channel:

- **Webhook URL** = `https://xxxx-xxx-xxx.ngrok-free.app/webhook`
- เปิด **Use webhook**
- กด **Verify** เพื่อทดสอบ (backend ต้องรันอยู่) — ควรได้ status 200

> โปรดักชัน: ชี้ Webhook URL ไปที่โดเมนจริงที่มี TLS แล้ว proxy เข้ามาที่ port 8000 (เช่นผ่าน Nginx/Caddy)

### 4. Webhook events ที่รองรับ

| event type | พฤติกรรม |
|---|---|
| `message` (text) | mark as read → embed → ค้น Weaviate → สร้างคำตอบด้วย Gemini → reply |
| `postback` | ถ้า `postback.data == "ปุ่มAนะ"` จะตอบข้อความ `"ปุ่มAนะ"` กลับไป (ตั้ง data นี้ได้จาก Rich menu / Flex button) |

LINE API ที่ backend เรียกออก:

- `POST https://api.line.me/v2/bot/message/reply` — ตอบข้อความ
- `POST https://api.line.me/v2/bot/chat/markAsRead` — ทำเครื่องหมายอ่านแล้ว

รูปแบบ payload ของ webhook อ้างอิง: <https://developers.line.biz/en/reference/messaging-api/#webhooks>

### 5. ตั้งค่า LIFF (สำหรับหน้า frontend)

1. LINE Developers Console → channel → แท็บ **LIFF** → **Add**
2. **Endpoint URL** = URL ที่ deploy หน้า `frontend/` (ตอน dev ใช้ ngrok ครอบ port 8888)
3. **Size** ตามต้องการ, **Scope** เลือก `profile`
4. นำ **LIFF ID** ที่ได้ไปแก้ใน `frontend/index.html` ที่บรรทัด `liff.init({liffId: "..."})`
   (ค่าปัจจุบัน: `2009838617-3r2S83yB`)

---

## Troubleshooting

| อาการ | สาเหตุที่พบบ่อย |
|---|---|
| `Weaviate Connection Error` / `เชื่อมต่อ Weaviate ไม่สำเร็จ` | ยังไม่ได้ `docker compose up -d` หรือ port 8080 ถูกใช้อยู่ |
| บอทตอบ `Error: System crash at response` | `GEMINI_API_KEY` ผิด/หมดโควตา หรือเรียก model ไม่ได้ |
| บอทตอบแต่ "ไม่มีข้อมูล" | ยังไม่ได้รัน `go run ./test/ingest/` หรือ class `Document` ว่าง |
| LINE **Verify** ไม่ผ่าน | backend ไม่ได้รัน, ngrok ปิด, หรือ Webhook URL ไม่ได้ลงท้าย `/webhook` |
| ได้ข้อความซ้ำ/ตอบชน | ยังไม่ได้ปิด Auto-reply / Greeting ใน LINE OA Manager |
| `Error: No .env file found` | ไม่ได้รันคำสั่งจากไดเรกทอรี `backend/` |
