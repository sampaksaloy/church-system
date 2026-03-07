package repository

import (
	"church-system/internal/models"
	"database/sql"
	"fmt"
)

type AnnouncementRepository struct {
	db *sql.DB
}

func NewAnnouncementRepository(db *sql.DB) *AnnouncementRepository {
	return &AnnouncementRepository{db: db}
}

func (r *AnnouncementRepository) FindAll(activeOnly bool) ([]*models.Announcement, error) {
	query := `
		SELECT a.id, a.title, a.content, a.category, a.is_pinned, a.is_active,
		       a.admin_id, COALESCE(ad.name, 'Admin') as admin_name,
		       a.created_at, a.updated_at
		FROM announcements a
		LEFT JOIN admins ad ON a.admin_id = ad.id`
	if activeOnly {
		query += ` WHERE a.is_active = true`
	}
	query += ` ORDER BY a.is_pinned DESC, a.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("find all announcements: %w", err)
	}
	defer rows.Close()

	var announcements []*models.Announcement
	for rows.Next() {
		a := &models.Announcement{}
		err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.IsPinned,
			&a.IsActive, &a.AdminID, &a.AdminName, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, a)
	}
	return announcements, nil
}

func (r *AnnouncementRepository) FindByID(id int) (*models.Announcement, error) {
	a := &models.Announcement{}
	err := r.db.QueryRow(`
		SELECT a.id, a.title, a.content, a.category, a.is_pinned, a.is_active,
		       a.admin_id, COALESCE(ad.name, 'Admin') as admin_name,
		       a.created_at, a.updated_at
		FROM announcements a
		LEFT JOIN admins ad ON a.admin_id = ad.id
		WHERE a.id = $1`, id,
	).Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.IsPinned,
		&a.IsActive, &a.AdminID, &a.AdminName, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find announcement by id: %w", err)
	}
	return a, nil
}

func (r *AnnouncementRepository) Create(title, content, category string, isPinned bool, adminID int) error {
	_, err := r.db.Exec(
		`INSERT INTO announcements (title, content, category, is_pinned, admin_id) VALUES ($1, $2, $3, $4, $5)`,
		title, content, category, isPinned, adminID,
	)
	return err
}

func (r *AnnouncementRepository) Update(id int, title, content, category string, isPinned, isActive bool) error {
	_, err := r.db.Exec(
		`UPDATE announcements SET title=$1, content=$2, category=$3, is_pinned=$4, is_active=$5, updated_at=CURRENT_TIMESTAMP WHERE id=$6`,
		title, content, category, isPinned, isActive, id,
	)
	return err
}

func (r *AnnouncementRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM announcements WHERE id = $1`, id)
	return err
}

func (r *AnnouncementRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM announcements WHERE is_active = true`).Scan(&count)
	return count, err
}

func (r *AnnouncementRepository) FindLatest(limit int) ([]*models.Announcement, error) {
	rows, err := r.db.Query(`
		SELECT id, title, content, category, is_pinned, is_active, admin_id,
		       COALESCE((SELECT name FROM admins WHERE id = admin_id), 'Admin'),
		       created_at, updated_at
		FROM announcements
		WHERE is_active = true
		ORDER BY is_pinned DESC, created_at DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Announcement
	for rows.Next() {
		a := &models.Announcement{}
		rows.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.IsPinned,
			&a.IsActive, &a.AdminID, &a.AdminName, &a.CreatedAt, &a.UpdatedAt)
		list = append(list, a)
	}
	return list, nil
}
