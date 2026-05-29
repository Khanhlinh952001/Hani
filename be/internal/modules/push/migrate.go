package push

import "be/internal/db"

// AutoMigrate creates push notification tables.
func AutoMigrate() error {
	return db.DB.AutoMigrate(&UserDevice{}, &Notification{})
}
