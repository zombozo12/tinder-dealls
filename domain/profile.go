package domain

import "time"

type Profile struct {
	ID         int64      `gorm:"primaryKey" json:"id" faker:"-"`
	UserID     int64      `gorm:"not null" json:"user_id" faker:"-"`
	Name       string     `gorm:"not null" json:"name" faker:"name"`
	Pic        string     `json:"pic" faker:"url"`
	Gender     string     `json:"gender" faker:"-"`
	InterestIn string     `json:"interest_in" faker:"-"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP()" json:"created_at" faker:"-"`
	UpdatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP()" json:"updated_at" faker:"-"`
	DeletedAt  *time.Time `gorm:"default:null" json:"deleted_at,omitempty" faker:"-"`
}

type ProfileRequest struct {
	Name       string `json:"name" validate:"required"`
	Gender     string `json:"gender" validate:"required"`
	InterestIn string `json:"interest_in" validate:"required"`
}

type UpdateProfilePicRequest struct {
	Pic string `json:"pic" validate:"required"`
}

type ProfileResponse struct {
	UserID     int64  `json:"user_id"`
	Name       string `json:"name"`
	Pic        string `json:"pic"`
	Gender     string `json:"gender"`
	InterestIn string `json:"interest_in"`
}
