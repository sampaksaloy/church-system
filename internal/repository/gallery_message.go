package repository

import (
	"church-system/internal/models"
	"database/sql"
	"fmt"
)

type GalleryRepository struct {
	db *sql.DB
}

func NewGalleryRepository(db *sql.DB) *GalleryRepository {
	return &GalleryRepository{db: db}
}

func (r *GalleryRepository) FindAll(activeOnly bool) ([]*models.GalleryPhoto, error) {
	query := `
		SELECT g.id, g.title, g.description, g.filename, g.filepath,
		       g.category, g.is_active, g.admin_id,
		       COALESCE(a.name, 'Admin'), g.created_at
		FROM gallery g
		LEFT JOIN admins a ON g.admin_id = a.id`
	if activeOnly {
		query += ` WHERE g.is_active = true`
	}
	query += ` ORDER BY g.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("find all gallery: %w", err)
	}
	defer rows.Close()

	var photos []*models.GalleryPhoto
	for rows.Next() {
		p := &models.GalleryPhoto{}
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Filename, &p.Filepath,
			&p.Category, &p.IsActive, &p.AdminID, &p.AdminName, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, nil
}

func (r *GalleryRepository) FindByID(id int) (*models.GalleryPhoto, error) {
	p := &models.GalleryPhoto{}
	err := r.db.QueryRow(`
		SELECT g.id, g.title, g.description, g.filename, g.filepath,
		       g.category, g.is_active, g.admin_id,
		       COALESCE(a.name, 'Admin'), g.created_at
		FROM gallery g
		LEFT JOIN admins a ON g.admin_id = a.id
		WHERE g.id = $1`, id,
	).Scan(&p.ID, &p.Title, &p.Description, &p.Filename, &p.Filepath,
		&p.Category, &p.IsActive, &p.AdminID, &p.AdminName, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find gallery by id: %w", err)
	}
	return p, nil
}

func (r *GalleryRepository) Create(title, description, filename, filepath, category string, adminID int) error {
	_, err := r.db.Exec(
		`INSERT INTO gallery (title, description, filename, filepath, category, admin_id) VALUES ($1, $2, $3, $4, $5, $6)`,
		title, description, filename, filepath, category, adminID,
	)
	return err
}

func (r *GalleryRepository) UpdateActive(id int, isActive bool) error {
	_, err := r.db.Exec(`UPDATE gallery SET is_active=$1 WHERE id=$2`, isActive, id)
	return err
}

func (r *GalleryRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM gallery WHERE id = $1`, id)
	return err
}

func (r *GalleryRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM gallery WHERE is_active = true`).Scan(&count)
	return count, err
}

func (r *GalleryRepository) FindLatest(limit int) ([]*models.GalleryPhoto, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, filename, filepath, category, is_active, admin_id,
		       COALESCE((SELECT name FROM admins WHERE id = admin_id), 'Admin'), created_at
		FROM gallery
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.GalleryPhoto
	for rows.Next() {
		p := &models.GalleryPhoto{}
		rows.Scan(&p.ID, &p.Title, &p.Description, &p.Filename, &p.Filepath,
			&p.Category, &p.IsActive, &p.AdminID, &p.AdminName, &p.CreatedAt)
		list = append(list, p)
	}
	return list, nil
}

// ─── Messages ────────────────────────────────────────────────────────────────

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) FindAll() ([]*models.Message, error) {
	rows, err := r.db.Query(`
		SELECT id, sender_name, email, phone, subject, message, is_read, created_at
		FROM messages ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("find all messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		m := &models.Message{}
		err := rows.Scan(&m.ID, &m.SenderName, &m.Email, &m.Phone,
			&m.Subject, &m.Message, &m.IsRead, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *MessageRepository) FindByID(id int) (*models.Message, error) {
	m := &models.Message{}
	err := r.db.QueryRow(`
		SELECT id, sender_name, email, phone, subject, message, is_read, created_at
		FROM messages WHERE id = $1`, id,
	).Scan(&m.ID, &m.SenderName, &m.Email, &m.Phone, &m.Subject, &m.Message, &m.IsRead, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MessageRepository) Create(senderName, email, phone, subject, message string) error {
	_, err := r.db.Exec(
		`INSERT INTO messages (sender_name, email, phone, subject, message) VALUES ($1, $2, $3, $4, $5)`,
		senderName, email, phone, subject, message,
	)
	return err
}

func (r *MessageRepository) MarkRead(id int) error {
	_, err := r.db.Exec(`UPDATE messages SET is_read=true WHERE id=$1`, id)
	return err
}

func (r *MessageRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM messages WHERE id=$1`, id)
	return err
}

func (r *MessageRepository) CountUnread() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM messages WHERE is_read = false`).Scan(&count)
	return count, err
}
