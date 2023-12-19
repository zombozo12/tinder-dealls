package domain

import "time"

type Matched struct {
	ID        int64      `json:"id"`
	UserAID   int64      `json:"user_a_id"`
	UserBID   int64      `json:"user_b_id"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP()" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:null" json:"deleted_at,omitempty"`
}

type MatchRequest struct {
	UserID       int64 `json:"user_id"`
	TargetUserID int64 `json:"target_user_id"`
}

type MatchResponse struct {
	ID           int64 `json:"id"`
	UserID       int64 `json:"user_id"`
	TargetUserID int64 `json:"target_user_id"`
}
