# Announcer - Backend

## Installing
ติดตั้ง <a href="https://golang.org/dl/">Go</a> บนเครื่องที่จะทำการรันตัวโปรเจค

## Getting Started
รันตัวโปรเจค
```bash
go run server.go
```
โปรเจคจะรันที่ http://localhost:8000

## Learn More
สามารถแก้ไขการตั้งค่าต่าง ๆ ได้ที่ไฟล์ .env
1. REDIRECT_URI คือ Path สำหรับ Line Login เพื่อใช้ Redirect ไปยังหน้า Login ของโปรเจค ในที่นี้คือ http://localhost:3000/login?social=line
```bash
REDIRECT_URI=http://localhost:3000/login?social=line
```
