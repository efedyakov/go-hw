package storage

import "time"

type Event struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StartAt     time.Time `db:"start_at"`
	EndAt       time.Time `db:"end_at"`
	NotifyAt    time.Time `db:"notify_at"`
}
