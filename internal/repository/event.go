package repository

import (
	"church-system/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) FindAll(activeOnly bool) ([]*models.Event, error) {
	query := `
		SELECT e.id, e.title, e.description, e.location, e.event_date,
		       COALESCE(e.start_time::text, ''), COALESCE(e.end_time::text, ''),
		       e.category, e.is_recurring, e.is_active, e.admin_id,
		       COALESCE(a.name, 'Admin'), e.created_at, e.updated_at
		FROM events e
		LEFT JOIN admins a ON e.admin_id = a.id`
	if activeOnly {
		query += ` WHERE e.is_active = true`
	}
	query += ` ORDER BY e.event_date ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("find all events: %w", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		e := &models.Event{}
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.EventDate,
			&e.StartTime, &e.EndTime, &e.Category, &e.IsRecurring, &e.IsActive,
			&e.AdminID, &e.AdminName, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *EventRepository) FindByID(id int) (*models.Event, error) {
	e := &models.Event{}
	err := r.db.QueryRow(`
		SELECT e.id, e.title, e.description, e.location, e.event_date,
		       COALESCE(e.start_time::text, ''), COALESCE(e.end_time::text, ''),
		       e.category, e.is_recurring, e.is_active, e.admin_id,
		       COALESCE(a.name, 'Admin'), e.created_at, e.updated_at
		FROM events e
		LEFT JOIN admins a ON e.admin_id = a.id
		WHERE e.id = $1`, id,
	).Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.EventDate,
		&e.StartTime, &e.EndTime, &e.Category, &e.IsRecurring, &e.IsActive,
		&e.AdminID, &e.AdminName, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find event by id: %w", err)
	}
	return e, nil
}

func (r *EventRepository) Create(title, description, location, category string, eventDate time.Time, startTime, endTime string, isRecurring bool, adminID int) error {
	var startT, endT interface{}
	if startTime != "" {
		startT = startTime
	}
	if endTime != "" {
		endT = endTime
	}
	_, err := r.db.Exec(
		`INSERT INTO events (title, description, location, event_date, start_time, end_time, category, is_recurring, admin_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		title, description, location, eventDate, startT, endT, category, isRecurring, adminID,
	)
	return err
}

func (r *EventRepository) Update(id int, title, description, location, category string, eventDate time.Time, startTime, endTime string, isRecurring, isActive bool) error {
	var startT, endT interface{}
	if startTime != "" {
		startT = startTime
	}
	if endTime != "" {
		endT = endTime
	}
	_, err := r.db.Exec(
		`UPDATE events SET title=$1, description=$2, location=$3, event_date=$4,
		 start_time=$5, end_time=$6, category=$7, is_recurring=$8, is_active=$9,
		 updated_at=CURRENT_TIMESTAMP WHERE id=$10`,
		title, description, location, eventDate, startT, endT, category, isRecurring, isActive, id,
	)
	return err
}

func (r *EventRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM events WHERE id = $1`, id)
	return err
}

func (r *EventRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM events WHERE is_active = true`).Scan(&count)
	return count, err
}

func (r *EventRepository) FindUpcoming(limit int) ([]*models.Event, error) {
	rows, err := r.db.Query(`
		SELECT e.id, e.title, e.description, e.location, e.event_date,
		       COALESCE(e.start_time::text, ''), COALESCE(e.end_time::text, ''),
		       e.category, e.is_recurring, e.is_active, e.admin_id,
		       COALESCE(a.name, 'Admin'), e.created_at, e.updated_at
		FROM events e
		LEFT JOIN admins a ON e.admin_id = a.id
		WHERE e.is_active = true AND e.event_date >= CURRENT_DATE
		ORDER BY e.event_date ASC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Event
	for rows.Next() {
		e := &models.Event{}
		rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.EventDate,
			&e.StartTime, &e.EndTime, &e.Category, &e.IsRecurring, &e.IsActive,
			&e.AdminID, &e.AdminName, &e.CreatedAt, &e.UpdatedAt)
		list = append(list, e)
	}
	return list, nil
}

func (r *EventRepository) FindByMonth(year, month int) ([]*models.Event, error) {
	rows, err := r.db.Query(`
		SELECT e.id, e.title, e.description, e.location, e.event_date,
		       COALESCE(e.start_time::text, ''), COALESCE(e.end_time::text, ''),
		       e.category, e.is_recurring, e.is_active, e.admin_id,
		       COALESCE(a.name, 'Admin'), e.created_at, e.updated_at
		FROM events e
		LEFT JOIN admins a ON e.admin_id = a.id
		WHERE e.is_active = true
		  AND EXTRACT(YEAR FROM e.event_date) = $1
		  AND EXTRACT(MONTH FROM e.event_date) = $2
		ORDER BY e.event_date ASC`, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Event
	for rows.Next() {
		e := &models.Event{}
		rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.EventDate,
			&e.StartTime, &e.EndTime, &e.Category, &e.IsRecurring, &e.IsActive,
			&e.AdminID, &e.AdminName, &e.CreatedAt, &e.UpdatedAt)
		list = append(list, e)
	}
	return list, nil
}
