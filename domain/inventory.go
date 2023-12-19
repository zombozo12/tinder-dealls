package domain

import "time"

type Inventory struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Swipes     int64     `json:"swipes"`
	Likes      int64     `json:"likes"`
	SuperLikes int64     `json:"super_likes"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP()" json:"updated_at"`
}

type CreateInventoryRequest struct {
	UserID     int64 `json:"user_id"`
	Likes      int64 `json:"likes"`
	SuperLikes int64 `json:"super_likes"`
}

type InventoryRequest struct {
	UserID int64 `json:"user_id"`
}

type InventoryResponse struct {
	UserID     int64 `json:"user_id"`
	Likes      int64 `json:"likes"`
	SuperLikes int64 `json:"super_likes"`
}
