package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("✅ Database connected successfully")
	return db, nil
}

func Migrate(db *sql.DB) error {
	migrations := []string{
		createAdminsTable,
		createAnnouncementsTable,
		createEventsTable,
		createGalleryTable,
		createMessagesTable,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	log.Println("✅ Database migrations applied successfully")
	return nil
}

const createAdminsTable = `
CREATE TABLE IF NOT EXISTS admins (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(150) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,
    role        VARCHAR(50) DEFAULT 'admin',
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createAnnouncementsTable = `
CREATE TABLE IF NOT EXISTS announcements (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    content     TEXT NOT NULL,
    category    VARCHAR(100) DEFAULT 'General',
    is_pinned   BOOLEAN DEFAULT FALSE,
    is_active   BOOLEAN DEFAULT TRUE,
    admin_id    INTEGER REFERENCES admins(id) ON DELETE SET NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createEventsTable = `
CREATE TABLE IF NOT EXISTS events (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    location     VARCHAR(255),
    event_date   DATE NOT NULL,
    start_time   TIME,
    end_time     TIME,
    category     VARCHAR(100) DEFAULT 'General',
    is_recurring BOOLEAN DEFAULT FALSE,
    is_active    BOOLEAN DEFAULT TRUE,
    admin_id     INTEGER REFERENCES admins(id) ON DELETE SET NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createGalleryTable = `
CREATE TABLE IF NOT EXISTS gallery (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    filename    VARCHAR(255) NOT NULL,
    filepath    VARCHAR(500) NOT NULL,
    category    VARCHAR(100) DEFAULT 'General',
    is_active   BOOLEAN DEFAULT TRUE,
    admin_id    INTEGER REFERENCES admins(id) ON DELETE SET NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createMessagesTable = `
CREATE TABLE IF NOT EXISTS messages (
    id          SERIAL PRIMARY KEY,
    sender_name VARCHAR(100) NOT NULL,
    email       VARCHAR(150) NOT NULL,
    phone       VARCHAR(30),
    subject     VARCHAR(255) NOT NULL,
    message     TEXT NOT NULL,
    is_read     BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
