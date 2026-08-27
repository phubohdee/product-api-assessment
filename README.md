# Product API Assessment

## 🚀 ขั้นตอนการติดตั้งและเริ่มต้นใช้งาน (Getting Started)

### 1. Clone หรือ Pull โปรเจกต์

```bash
# กรณี Clone ใหม่
git clone https://github.com/phubohdee/product-api-assessment.git
cd product-api-assessment

# หรือกรณีมีโปรเจกต์เดิมอยู่แล้ว
git pull origin main
```

---

### 2. คัดลอกและเปลี่ยนชื่อไฟล์ .env.example เป็น .env

```bash
cp .env.example .env
```

---

### 3. เปิดใช้งาน PostgreSQL Database (Docker)

```bash
docker-compose up -d
```

---

### 4. ติดตั้ง Go Dependencies

```bash
go mod tidy
```

---

### 5. รัน Database Migration (สร้างตารางใน Database)

```bash
make migrate-up
```

---

### 6. เริ่มต้นรัน API Server

```bash
make run
```

---

### 7. ทดสอบการใช้งานผ่าน Swagger UI

เมื่อเซิร์ฟเวอร์เริ่มทำงานแล้ว สามารถเปิดเบราว์เซอร์ไปที่:

```text
http://localhost:8080/v1/api-docs/index.html
```

---

## 🛠️ รายการคำสั่ง Makefile ทั้งหมด

| คำสั่ง | รายละเอียด |
|--------|------------|
| `make run` | รัน API Server |
| `make migrate-up` | สร้าง/อัปเดตตารางใน Database |
| `make migrate-down` | ย้อนกลับ Migration ใน Database |
| `make test` | รัน Unit Tests |
| `make test-integration` | รัน Unit Tests และ Integration Tests ทั้งหมด |
| `make gen-swagger` | สร้างเอกสาร Swagger UI ใหม่จาก Annotation |

---

## 🧪 การรัน Test

```bash
# รัน Unit Tests เท่านั้น (ไม่จำเป็นต้องเปิด DB)
make test

# รัน Integration Tests (ต้องเปิด PostgreSQL ใน Docker)
make test-integration
```

---

## ⚠️ Error Codes Reference

| Error Code | HTTP Status | คำอธิบาย |
|------------|-------------|----------|
| `INVALID_REQUEST` | 400 Bad Request | รูปแบบ JSON หรือ Request Body ไม่ถูกต้อง |
| `INVALID_NAME` | 400 Bad Request | ชื่อสินค้าว่างเปล่า |
| `INVALID_PRICE` | 400 Bad Request | ราคาหลักน้อยกว่าหรือเท่ากับ 0 |
| `INVALID_SALE_PRICE` | 400 Bad Request | ราคาลดพิเศษมากกว่าหรือเท่ากับราคาหลัก |
| `PRODUCT_NOT_FOUND` | 404 Not Found | ไม่พบสินค้าตาม ID ที่ระบุ |
| `INTERNAL_SERVER_ERROR` | 500 Internal Server Error | เกิดข้อผิดพลาดฝั่ง Database หรือ Server |
