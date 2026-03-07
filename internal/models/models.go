package models

import "time"

type Admin struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Announcement struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	IsPinned  bool      `json:"is_pinned"`
	IsActive  bool      `json:"is_active"`
	AdminID   *int      `json:"admin_id"`
	AdminName string    `json:"admin_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	EventDate   time.Time `json:"event_date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Category    string    `json:"category"`
	IsRecurring bool      `json:"is_recurring"`
	IsActive    bool      `json:"is_active"`
	AdminID     *int      `json:"admin_id"`
	AdminName   string    `json:"admin_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GalleryPhoto struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Filename    string    `json:"filename"`
	Filepath    string    `json:"filepath"`
	Category    string    `json:"category"`
	IsActive    bool      `json:"is_active"`
	AdminID     *int      `json:"admin_id"`
	AdminName   string    `json:"admin_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type Message struct {
	ID         int       `json:"id"`
	SenderName string    `json:"sender_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Subject    string    `json:"subject"`
	Message    string    `json:"message"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

// Form structs
type LoginForm struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type AnnouncementForm struct {
	Title    string `form:"title" binding:"required"`
	Content  string `form:"content" binding:"required"`
	Category string `form:"category"`
	IsPinned bool   `form:"is_pinned"`
	IsActive bool   `form:"is_active"`
}

type EventForm struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	Location    string `form:"location"`
	EventDate   string `form:"event_date" binding:"required"`
	StartTime   string `form:"start_time"`
	EndTime     string `form:"end_time"`
	Category    string `form:"category"`
	IsRecurring bool   `form:"is_recurring"`
	IsActive    bool   `form:"is_active"`
}

type ContactForm struct {
	SenderName string `form:"sender_name" binding:"required"`
	Email      string `form:"email" binding:"required,email"`
	Phone      string `form:"phone"`
	Subject    string `form:"subject" binding:"required"`
	Message    string `form:"message" binding:"required"`
}

type AdminProfileForm struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password"`
}
