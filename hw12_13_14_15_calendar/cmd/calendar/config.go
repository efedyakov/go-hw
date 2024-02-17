package main

import (
	"fmt"
	"time"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	App        AppConf
	Server     ServerConf
	Logger     LoggerConf
	PostgreSQL PostgreSQLConf
	// TODO
}

type LoggerConf struct {
	Level string `envconfig:"LOGGER_LEVEL" default:"info" required:"true"`
}

type AppConf struct {
	Name        string `envconfig:"APP_NAME" default:"calendar" required:"true"`
	StorageType string `envconfig:"APP_STORAGE_TYPE" required:"true"`
}

type ServerConf struct {
	Host         string        `envconfig:"SERVER_HOST" default:"127.0.0.1" required:"true"`
	Port         string        `envconfig:"SERVER_PORT" default:"80" required:"true"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"30s" required:"true"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"30s" required:"true"`
	IdleTimeout  time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"30s" required:"true"`
}

type PostgreSQLConf struct {
	Host     string `envconfig:"PG_HOST" required:"true"`
	Port     int    `envconfig:"PG_PORT" required:"true"`
	Username string `envconfig:"PG_USERNAME" required:"true"`
	Password string `envconfig:"PG_PASSWORD" required:"true"`
	Database string `envconfig:"PG_DATABASE" required:"true"`
}

func NewConfig() Config {
	return Config{}
}

func (p *PostgreSQLConf) BuildDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", p.Username, p.Password, p.Host, p.Port, p.Database)
}
