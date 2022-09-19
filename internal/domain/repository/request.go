package repository

import (
	"context"
	"github.com/pkg/errors"
)

var (
	ErrNotFount = errors.New("не удалось найти запись")
)

// Интерфейс репозитория запросов и ответов
type RequestRepositoryInterface interface {
	// Создает запись
	Create(ctx context.Context, row Row) (*Row, error)

	// Возвращет одну запись по id
	GetOneByID(ctx context.Context, id string) (*Row, error)

	// Возвращет все записи
	GetAll(ctx context.Context) ([]*Row, error)

	// Удаляет одну запись по id
	Delete(ctx context.Context, id string) error
}

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type Response struct {
	ID      string            `json:"id"`
	Status  string            `json:"status"`
	Headers map[string]string `json:"headers"`
	Length  int64             `json:"length"`
}

type Row struct {
	ID       string
	Request  Request
	Response Response
}
