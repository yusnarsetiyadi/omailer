# OMailer

> **O**pen **Mailer** — Multi-purpose notification bridge for email.
<!-- > **O**rganization **Mailer** — Multi-purpose notification bridge for email & WhatsApp. -->

[![Go Version](https://img.shields.io/badge/Go-1.25.0-00ADD8?logo=go)](https://go.dev/)
[![Echo v4](https://img.shields.io/badge/Echo-v4-4DB33D?logo=go)](https://echo.labstack.com/)
[![License](https://img.shields.io/badge/license-MIT-green)](#)
[![Docker](https://img.shields.io/badge/docker-ready-2496ED?logo=docker)](https://www.docker.com/)

---

## Features

| | Feature | Description |
|---|---|---|
| :incoming_envelope: | **Email API** | Send rich HTML emails with file attachments via REST API (bring-your-own SMTP). |
| :globe_with_meridians: | **Swagger UI** | Interactive API documentation at `/swagger/*`. |
| :whale: | **Dockerized** | Multi-stage Docker build with CI/CD via GitHub Actions. |
<!-- | :speech_balloon: | **WhatsApp Bot** | Multi-device WhatsApp client with group mention support and session persistence. |
| :alarm_clock: | **Scheduler** | Cron-based automated messaging (e.g., daily attendance reminders). | -->

---

## Tech Stack

| Category | Technology |
|---|---|
| **Language** | Go 1.25.0 |
| **HTTP Framework** | [Echo v4](https://echo.labstack.com/) |
| **Email** | [gomail v2](https://gopkg.in/gomail.v2) |
| **Database** | SQLite (WhatsApp session storage) |
| **Validation** | [go-playground/validator](https://github.com/go-playground/validator) |
| **Logging** | [Logrus](https://github.com/sirupsen/logrus) |
| **Docs** | Swagger / OpenAPI |
<!-- | **WhatsApp** | [WhatsMeow](https://go.mau.fi/whatsmeow) |
| **Scheduler** | [GoCron v2](https://github.com/go-co-op/gocron) | -->

---

## Quick Start

### Prerequisites
- Go 1.25+
- Docker (optional)

### Run Locally

```bash
git clone https://github.com/your-username/omailer.git
cd omailer
go run main.go
```

<!-- On first run, a **QR code** will appear in the terminal — scan it with WhatsApp (Linked Devices) to pair the bot. -->

### Run with Docker

```bash
docker build -t omailer .
docker run -p 9000:9000 omailer
```

---

## API Endpoints

### `POST /send` — Send email with attachments (multipart/form-data)

```bash
curl -X POST https://yusnar.my.id/omailer/send \
  -F "smtp_host=smtp.gmail.com" \
  -F "smtp_port=587" \
  -F "auth_email=you@gmail.com" \
  -F "auth_password=your-app-password" \
  -F "sender_name=John Doe" \
  -F "recipient=friend@example.com" \
  -F "subject=Hello from OMailer" \
  -F "body_html=<h1>Hi there!</h1><p>This is an OMailer test.</p>" \
  -F "files=@document.pdf"
```

**Parameters:**
| Field | Type | Required | Description |
|---|---|---|---|
| `smtp_host` | string | yes | SMTP server hostname |
| `smtp_port` | int | yes | SMTP server port |
| `auth_email` | string | yes | SMTP authentication email |
| `auth_password` | string | yes | SMTP authentication password |
| `sender_name` | string | yes | Display name of the sender |
| `recipient` | string | yes | Recipient email address |
| `subject` | string | yes | Email subject |
| `body_html` | string | yes | HTML email body |
| `files` | file | no | File attachment(s) |

### `GET /send/just-message` — Send email via URL-encoded JSON

```bash
curl "https://yusnar.my.id/omailer/send/just-message?data=%7B%22smtp_host%22%3A%22smtp.gmail.com%22%2C%22...%22%7D"
```

The `data` query parameter must be a URL-encoded JSON string with the same fields as above.

### `GET /` — Health check

```json
{
  "message": "Welcome to OMailer"
}
```

### `GET /swagger/*` — Swagger UI

Interactive API documentation served via Swagger.

---

<!-- ## WhatsApp Bot

OMailer runs a WhatsApp multi-device client that:

- **Pairs via QR code** on first launch (session persisted to SQLite).
- **Sends text messages** to individual numbers or WhatsApp groups.
- **Mentions all group members** when sending to a group.
- Resolves group names automatically (fuzzy matching).

### Scheduled Jobs

| Time (WIB) | Day | Message |
|---|---|---|
| 08:00 | Mon–Fri | Attendance check-in reminder |
| 17:00 | Mon–Fri | Attendance check-out reminder |

--- -->

## Project Structure

```
├── main.go                     # Entry point
├── Dockerfile                  # Multi-stage Docker build
├── go.mod / go.sum             # Go module dependencies
├── .github/workflows/main.yml  # CI/CD pipeline
│
├── internal/
│   ├── abstraction/            # Custom Echo context
│   ├── app/omailer/            # Handler, routes, and service layer
│   ├── dto/                    # Request DTOs
│   ├── http/                   # HTTP server setup & route registration
│   ├── middleware/              # CORS, logger, recover, error handler
│   └── scheduler/              # GoCron jobs & helpers
│
├── pkg/
│   ├── constant/               # App constants (port, version, env)
│   ├── general/                # HTML-to-plaintext converter, time helpers
│   ├── gomail/                 # SMTP email sender
│   ├── log/                    # Logrus initialization
│   ├── util/                   # Response builders & validator
│   └── whatsapp/               # WhatsApp client & helpers
│
└── docs/                       # Swagger specifications
```

---

## CI/CD

On every push to `main`, GitHub Actions automatically:

1. :hammer: Builds the Docker image
2. :rocket: Pushes to Docker Hub (`yusnars/omailer-personal-image:latest`)
3. :satellite: Deploys to VPS via SSH

---

## Environment

Configuration lives in `pkg/constant/constant.go`:

| Variable | Default | Description |
|---|---|---|
| `PORT` | `9000` | HTTP server port |
| `VERSION` | `1.0.0` | App version |
| `ENV` | `dev` | Environment name |

> SMTP credentials are passed **per request** — no server-side email config required.

---

## Security Notes

SMTP credentials are client-provided and never stored on the server.
<!-- - WhatsApp session is persisted locally in `data/session.db` (cleaned up on shutdown). -->

---

## License

MIT
