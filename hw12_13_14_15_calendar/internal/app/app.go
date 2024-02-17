package app

import (
	"context"
	"net/http"

	"github.com/efedyakov/go-hw/hw12_13_14_15_calendar/internal/storage"
)

// var (
//	errDateIsBooked = errors.New("this date is already booked")
// )

type App struct {
	name    string
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) (int, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID int) error
	ListEvent(ctx context.Context) ([]storage.Event, error)
	GetEvent(ctx context.Context, eventID int) (*storage.Event, error)
}

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func New(name string, logger Logger, storage Storage) *App {
	return &App{
		name:    name,
		logger:  logger,
		storage: storage,
	}
}

func (a *App) GetName() string {
	return a.name
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) (int, error) {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, eventID int) error {
	return a.storage.DeleteEvent(ctx, eventID)
}

func (a *App) ListEvent(ctx context.Context) ([]storage.Event, error) {
	return a.storage.ListEvent(ctx)
}

func (a *App) GetEvent(ctx context.Context, eventID int) (*storage.Event, error) {
	return a.storage.GetEvent(ctx, eventID)
}
