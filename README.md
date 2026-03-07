# ✝ Holy Trinity Church – Event & Announcement System

A full-featured church management website built with **Go + Gin + PostgreSQL**.

---

## Features

### Public (Members / Visitors)
- 🏠 **Home** – Welcome page with latest announcements, upcoming events, and photo previews
- 📋 **Announcements** – Browse and read all parish announcements, filterable by category
- 📅 **Event Calendar** – Monthly calendar view + upcoming events list
- 📷 **Photo Gallery** – Lightbox photo gallery, filterable by category
- ✉️ **Contact** – Send inquiries to the church office

### Admin Panel
- 🔐 **Secure Login** – Session-based authentication
- 📊 **Dashboard** – Overview statistics + recent activity
- 📋 **Announcements** – Create, edit, delete, pin announcements
- 📅 **Events** – Manage church events with date/time/location
- 📷 **Gallery** – Upload and manage photos
- ✉️ **Messages** – Read and reply to contact form submissions
- 👤 **Profile** – Update admin name, email, and password

---

## Tech Stack

| Layer    | Technology      |
|----------|----------------|
| Language | Go 1.21+        |
| Web Framework | Gin     |
| Database | PostgreSQL      |
| Sessions | gorilla/sessions via gin-contrib/sessions |
| Auth     | bcrypt          |
| Frontend | Pure HTML/CSS (no frontend build step) |

---

## Prerequisites

- Go 1.21+
- PostgreSQL 13+
- (Optional) `make`

---

## Setup Instructions

### 1. Clone / Unzip the project
```bash
unzip church-system.zip
cd church-system
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. Create PostgreSQL database
```bash
createdb church_db
# OR using psql:
psql -U postgres -c "CREATE DATABASE church_db;"
```

### 4. Configure environment variables
```bash
cp .env.example .env
# Edit .env with your database credentials
```

Your `.env` file should look like:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=church_db
DB_SSLMODE=disable
SERVER_PORT=8080
SESSION_SECRET=change-this-to-random-string
ADMIN_DEFAULT_EMAIL=admin@church.com
ADMIN_DEFAULT_PASSWORD=Admin@123
```

### 5. Run the application
```bash
go run ./cmd/server/main.go
# OR
make run
```

The app will:
- Connect to PostgreSQL
- Run database migrations automatically
- Create a default admin account (if none exists)
- Start the server

### 6. Access the system

| URL | Description |
|-----|-------------|
| http://localhost:8080 | Public website |
| http://localhost:8080/admin/login | Admin login |

**Default admin credentials:**
- Email: `admin@church.com`
- Password: `Admin@123`

> ⚠️ Change the password immediately after first login!

### 7. (Optional) Load sample data
```bash
make db-seed
# OR
psql -d church_db -f db/seed.sql
```

---

## Project Structure

```
church-system/
├── cmd/server/
│   └── main.go              # Entry point, router setup
├── internal/
│   ├── config/
│   │   └── config.go        # Environment config
│   ├── database/
│   │   └── database.go      # DB connection + migrations
│   ├── models/
│   │   └── models.go        # Data models & form structs
│   ├── repository/
│   │   ├── admin.go         # Admin CRUD
│   │   ├── announcement.go  # Announcement CRUD
│   │   ├── event.go         # Event CRUD
│   │   └── gallery_message.go # Gallery & Message CRUD
│   ├── handlers/
│   │   ├── public.go        # Public page handlers
│   │   └── admin.go         # Admin panel handlers
│   └── middleware/
│       └── auth.go          # Auth middleware
├── templates/
│   ├── partials/
│   │   └── base.html        # Public base layout
│   ├── public/              # Public page templates
│   │   ├── home.html
│   │   ├── announcements.html
│   │   ├── announcement_detail.html
│   │   ├── events.html
│   │   ├── gallery.html
│   │   └── contact.html
│   └── admin/               # Admin panel templates
│       ├── base.html        # Admin base layout
│       ├── login.html
│       ├── dashboard.html
│       ├── announcements.html
│       ├── announcement_form.html
│       ├── events.html
│       ├── event_form.html
│       ├── gallery.html
│       ├── gallery_form.html
│       ├── messages.html
│       ├── message_view.html
│       └── profile.html
├── static/
│   └── images/uploads/      # Uploaded photos stored here
├── db/
│   └── seed.sql             # Sample data
├── .env.example
├── go.mod
├── Makefile
└── README.md
```

---

## Database Schema

```sql
admins          -- Church staff accounts
announcements   -- Parish news and notices
events          -- Church activities and schedules
gallery         -- Photo uploads
messages        -- Contact form submissions
```

---

## Building for Production

```bash
make build
# Binary will be at ./bin/server
./bin/server
```

Set `GIN_MODE=release` in production:
```bash
GIN_MODE=release ./bin/server
```
