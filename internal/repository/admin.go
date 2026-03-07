package repository

import (
	"church-system/internal/models"
	"database/sql"
	"fmt"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) FindByEmail(email string) (*models.Admin, error) {
	admin := &models.Admin{}
	err := r.db.QueryRow(
		`SELECT id, name, email, password, role, created_at, updated_at FROM admins WHERE email = $1`,
		email,
	).Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Password, &admin.Role, &admin.CreatedAt, &admin.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find admin by email: %w", err)
	}
	return admin, nil
}

func (r *AdminRepository) FindByID(id int) (*models.Admin, error) {
	admin := &models.Admin{}
	err := r.db.QueryRow(
		`SELECT id, name, email, password, role, created_at, updated_at FROM admins WHERE id = $1`,
		id,
	).Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Password, &admin.Role, &admin.CreatedAt, &admin.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find admin by id: %w", err)
	}
	return admin, nil
}

func (r *AdminRepository) Create(name, email, hashedPassword string) error {
	_, err := r.db.Exec(
		`INSERT INTO admins (name, email, password) VALUES ($1, $2, $3)`,
		name, email, hashedPassword,
	)
	return err
}

func (r *AdminRepository) Update(id int, name, email, hashedPassword string) error {
	if hashedPassword != "" {
		_, err := r.db.Exec(
			`UPDATE admins SET name=$1, email=$2, password=$3, updated_at=CURRENT_TIMESTAMP WHERE id=$4`,
			name, email, hashedPassword, id,
		)
		return err
	}
	_, err := r.db.Exec(
		`UPDATE admins SET name=$1, email=$2, updated_at=CURRENT_TIMESTAMP WHERE id=$3`,
		name, email, id,
	)
	return err
}

func (r *AdminRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM admins`).Scan(&count)
	return count, err
}
