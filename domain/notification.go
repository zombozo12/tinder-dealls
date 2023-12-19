package domain

import "time"

type Notification struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Message   string     `json:"message"`
	IsRead    bool       `json:"is_read"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP()" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:null" json:"deleted_at,omitempty"`
}

type NotificationRequest struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type NotificationResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
