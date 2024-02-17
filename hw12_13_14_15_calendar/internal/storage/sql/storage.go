package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/efedyakov/go-hw/hw12_13_14_15_calendar/internal/storage"

	"github.com/jmoiron/sqlx"
)

var errTimeoutClosingConn = errors.New("timeout closing connection")

type Storage struct {
	db *sqlx.DB
}

func New(dsn string) (*Storage, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s Storage) CreateEvent(ctx context.Context, event storage.Event) (int, error) {
	var lastInsertID int

	args := map[string]interface{}{
		"user_id":     event.UserID,
		"title":       event.Title,
		"description": event.Description,
		"start_at":    event.StartAt,
		"end_at":      event.EndAt,
		"notify_at":   event.NotifyAt,
	}
	query := `
insert into events
(user_id, title, description, start_at, end_at, notify_at) values 
(:user_id, :title, :description, :start_at, :end_at, :notify_at)
RETURNING id;`

	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	if err := stmt.GetContext(ctx, &lastInsertID, args); err != nil {
		return 0, err
	}
	return lastInsertID, err
}

func (s Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	args := map[string]interface{}{
		"id":          event.ID,
		"title":       event.Title,
		"description": event.Description,
		"start_at":    event.StartAt,
		"end_at":      event.EndAt,
		"notify_at":   event.NotifyAt,
	}
	query := `
update events set
(title, description, start_at, end_at, notify_at) =
(:title, :description, :start_at, :end_at, :notify_at) 
where id = :id;`

	_, err := s.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) DeleteEvent(ctx context.Context, eventID int) error {
	args := map[string]interface{}{
		"id": eventID,
	}
	query := `delete from events where id = :id`

	rows, err := s.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (s Storage) ListEvent(ctx context.Context) ([]storage.Event, error) {
	query := `select id, user_id, title, description, start_at, end_at, notify_at from events;`
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]storage.Event, 0)
	for rows.Next() {
		var event storage.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (s Storage) GetEvent(ctx context.Context, eventID int) (*storage.Event, error) {
	var event storage.Event

	args := map[string]interface{}{
		"id": eventID,
	}
	query := `select id, user_id, title, description, start_at, end_at, notify_at from events where id = :id;`

	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &event, args)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (s Storage) Connect(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s Storage) Close(ctx context.Context) error {
	var (
		ch             = make(chan error)
		newCtx, cancel = context.WithTimeout(ctx, 3*time.Second)
	)
	defer cancel()

	go func() {
		ch <- s.db.Close()
	}()

	select {
	case <-newCtx.Done():
		return errTimeoutClosingConn
	case <-ch:
		return nil
	}
}
