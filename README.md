# 🎵 Jukebox Analytics Service

Centralized backend service for collecting and analyzing jukebox playback data.

The service exposes an HTTP API that allows:

* registering playback events
* updating track prices
* retrieving top tracks
* calculating revenue statistics

---

## 🚀 How to Run

### 🐳 Run with Docker (recommended)

Make sure you have Docker and Docker Compose installed.

```bash
docker compose up --build
```

This will start:

* PostgreSQL database
* Application server

Wait for this message to appear in logs

```
[timestamp] Analytics server is running...
```

### 🌐 Available Services

* API: http://localhost:8080
* Swagger UI: http://localhost:8081
* PostgreSQL: localhost:5432

---

## ⚙️ Run Manually

### 1. Set environment variable

```bash
export DATABASE_URL=postgres://jukebox:secret@localhost:5432/jukebox?sslmode=disable
```

### 2. Run migrations

```bash
migrate -path infrastructure/migrations -database "$DATABASE_URL" up
```

### 3. Start server

```bash
go run ./cmd
```

Server will be available at:

```
http://localhost:8080
```

---

## 📬 API Endpoints

### ➤ Register Playback

```
POST /api/v1/logs
```

Request:

```json
{
  "track_id": 1,
  "amount_paid": 150
}
```

---

### ➤ Update Track Price

```
PATCH /api/v1/tracks/{id}/price
```

Request:

```json
{
  "newPrice": 200
}
```

---

### ➤ Get Top Tracks

```
GET /api/v1/stats/top
```

---

### ➤ Get Revenue Statistics

```
GET /api/v1/stats/revenue
```

---

## 📖 API Documentation

Swagger/OpenAPI documentation is available at:

```
http://localhost:8080/swagger/index.html
```

---

## 📝 Notes

* If running manually, make sure:

  * `DATABASE_URL` is set
  * migrations from `infrastructure/migrations` are applied
* API uses JSON for all requests and responses
* All endpoints follow `/api/v1/...` prefix
