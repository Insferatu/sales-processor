# Sales Processor

A lightweight Go service that processes sales events and distributes them to Google Sheets and Telegram. Built as a cost-effective replacement for Zapier webhooks.

## Overview

This service exposes REST endpoints that receive sales event webhooks, then:
1. Append the sale data as a row to a Google Sheet
2. Send a notification message to a Telegram channel

Both operations run concurrently for optimal performance.

## Endpoints

### 3D Toy Sales
```
POST /3d-toy-sale/
```

**Request body:**
```json
{
    "time": "2025-11-18 22:45:48",
    "item": "Марк",
    "material": "Золотой",
    "price": "40",
    "paymentType": "Карта"
}
```

**Google Sheets columns:** Время продажи | Фигурка | Пластик | Цена | Тип оплаты

**Telegram message:**
```
Фигурка: Марк
Пластик: Золотой
Продано за: 40
Тип оплаты: Карта
```

### Jewelry Sales
```
POST /jewelry-sale/
```

**Request body:**
```json
{
    "time": "2025-11-18 22:45:48",
    "product": "Браслет",
    "price": "40",
    "paymentType": "Карта"
}
```

**Google Sheets columns:** Время продажи | Товар | Цена | Тип оплаты

**Telegram message:**
```
Товар: Браслет
Продано за: 40
Тип оплаты: Карта
```

### Health Check
```
GET /health
```

Returns `{"status": "healthy"}` when the service is running.

## Local Development

### Prerequisites

- Go 1.22+
- Google Cloud service account with Sheets API access
- Telegram bot token

### Environment Variables

Create a `.env` file:

```bash
# Google Sheets authentication
GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account.json"

# 3D Toy Sale spreadsheet
TOY_SPREADSHEET_ID="your-toy-spreadsheet-id"
TOY_SHEET_NAME="Sales"

# Jewelry Sale spreadsheet
JEWELRY_SPREADSHEET_ID="your-jewelry-spreadsheet-id"
JEWELRY_SHEET_NAME="Sales"

# Telegram
TELEGRAM_BOT_TOKEN="your-bot-token"
TELEGRAM_CHAT_ID="your-chat-id"
```

### Build

```bash
go build ./...
```

### Test

```bash
go test ./... -v
```

### Run

```bash
source .env && go run ./cmd/server
```

The server starts on port 8080 by default. Set the `PORT` environment variable to change it.

### Test Requests

```bash
# 3D Toy Sale
curl -X POST http://localhost:8080/3d-toy-sale/ \
  -H "Content-Type: application/json" \
  -d '{
    "time": "2025-11-18 22:45:48",
    "item": "Марк",
    "material": "Золотой",
    "price": "40",
    "paymentType": "Карта"
  }'

# Jewelry Sale
curl -X POST http://localhost:8080/jewelry-sale/ \
  -H "Content-Type: application/json" \
  -d '{
    "time": "2025-11-18 22:45:48",
    "product": "Браслет",
    "price": "40",
    "paymentType": "Карта"
  }'
```

## Project Structure

```
sales-processor/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── health.go            # Health check endpoint
│   │   ├── toy_sale.go          # 3D toy sale handler
│   │   ├── toy_sale_test.go
│   │   ├── jewelry_sale.go      # Jewelry sale handler
│   │   └── jewelry_sale_test.go
│   ├── sheets/
│   │   └── client.go            # Google Sheets API client
│   └── telegram/
│       └── client.go            # Telegram Bot API client
├── .github/
│   └── workflows/
│       └── ci.yml               # GitHub Actions CI
├── cloudbuild.yaml              # Google Cloud Build config
├── Dockerfile
├── go.mod
└── go.sum
```

## CI/CD

### Continuous Integration (GitHub Actions)

On every push and pull request to `main`:
- Runs `go vet` for static analysis
- Runs `staticcheck` for additional linting
- Executes all tests with race detection
- Builds the binary and Docker image

See `.github/workflows/ci.yml` for details.

### Continuous Deployment (Google Cloud Build)

On every push to `main`:
1. Builds the Docker image
2. Pushes to Google Artifact Registry
3. Deploys to Cloud Run with zero-downtime

The deployment is configured in `cloudbuild.yaml`.

### Cloud Run Configuration

- **Memory:** 128Mi
- **CPU:** 1
- **Min instances:** 0 (scale to zero when idle)
- **Max instances:** 10
- **Request timeout:** 30s
- **Region:** europe-west1

## Deployment

### Google Cloud Setup

1. **Enable APIs:**
   - Cloud Run API
   - Cloud Build API
   - Secret Manager API
   - Google Sheets API
   - Artifact Registry API

2. **Create Service Accounts:**
   - `sales-processor-sa` - Runtime service account for Cloud Run
   - Cloud Build service account needs: Cloud Run Admin, Artifact Registry Writer, Secret Manager Secret Accessor

3. **Store Secrets in Secret Manager:**
   - `telegram-bot-token`
   - `telegram-chat-id`
   - `toy-spreadsheet-id`
   - `toy-sheet-name`
   - `jewelry-spreadsheet-id`
   - `jewelry-sheet-name`

4. **Share Google Sheets** with the service account email.

5. **Connect GitHub to Cloud Build:**
   - Create a trigger on push to `main` branch
   - Use `cloudbuild.yaml` as the build configuration

### Manual Deployment

To trigger a deployment manually, push to the `main` branch or use Cloud Build console to retry a build.

### Stopping the Service

To stop the service and avoid charges:
1. Go to Cloud Run console
2. Select the service
3. Click "Edit & Deploy New Revision"
4. Set "Minimum instances" to 0
5. Optionally delete the service entirely

## Cost Optimization

- **Scale to zero:** No charges when there's no traffic
- **Minimal resources:** 128Mi memory, shared CPU
- **Artifact Registry:** Stores only latest images
- **Secret Manager:** Free tier covers typical usage

Estimated cost: ~$0-5/month with moderate traffic, $0 when idle.

## License

Private project.
