package users

import "time"

type User struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null"`
	Password    string    `json:"-" gorm:"not null"`
	PhoneNumber string    `json:"phone_number"`
	Provider    string    `json:"provider"`
	ProviderId  string    `json:"provider_id"`
	Avatar      string    `json:"avatar"`
	Level       int       `json:"level"` // topik 1 2 3 4 5 6
	Address     string    `json:"address"`
	Status      int       `json:"status"`
	Role        int       `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
