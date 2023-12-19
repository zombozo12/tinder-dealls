package domain

import "time"

type User struct {
	ID             int64      `gorm:"primaryKey" json:"id" faker:"-"`
	Email          string     `gorm:"unique;not null" json:"email" faker:"email"`
	Password       string     `gorm:"not null" json:"password,omitempty" faker:"-"`
	AccessToken    string     `json:"access_token,omitempty" faker:"-"`
	TokenExpiredAt *time.Time `json:"token_expired_at,omitempty" faker:"-"`
	CreatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP()" json:"created_at" faker:"-"`
	UpdatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP()" json:"updated_at" faker:"-"`
	DeletedAt      *time.Time `gorm:"default:null" json:"deleted_at,omitempty" faker:"-"`
}

type UpdateTokenRequest struct {
	AccessToken    string    `validate:"required"`
	TokenExpiredAt time.Time `validate:"required"`
	UpdatedAt      time.Time
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8,max=32"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=32"`
}
