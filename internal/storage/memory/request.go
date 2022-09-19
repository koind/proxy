package memory

import (
	"context"
	"github.com/koind/proxy/internal/domain/repository"
	"github.com/pkg/errors"
	"sync"
)

// Создает фиктивный репозиторий запросов и ответов
func NewRequestRepository() *RequestRepository {
	return &RequestRepository{
		DB: make(map[string]*repository.Row),
	}
}

// Фиктивный репозиторий запросов и ответов
type RequestRepository struct {
	sync.RWMutex
	DB map[string]*repository.Row
}

// Создает запись
func (r *RequestRepository) Create(ctx context.Context, row repository.Row) (*repository.Row, error) {
	_, err := r.GetOneByID(ctx, row.ID)
	if err != nil && !errors.Is(err, repository.ErrNotFount) {
		return nil, err
	}

	r.Lock()
	defer r.Unlock()

	r.DB[row.ID] = &row

	return &row, nil
}

// Возвращет одну запись по id
func (r *RequestRepository) GetOneByID(ctx context.Context, id string) (*repository.Row, error) {
	r.RLock()
	defer r.RUnlock()

	row, has := r.DB[id]
	if !has {
		return nil, repository.ErrNotFount
	}

	return row, nil
}

// Возвращет все записи
func (r *RequestRepository) GetAll(ctx context.Context) ([]*repository.Row, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.DB) <= 0 {
		return nil, nil
	}

	list := make([]*repository.Row, 0, len(r.DB))

	for _, Row := range r.DB {
		list = append(list, Row)
	}

	return list, nil
}

// Удаляет одну запись по id
func (r *RequestRepository) Delete(ctx context.Context, id string) error {
	_, err := r.GetOneByID(ctx, id)
	if err != nil {
		return err
	}

	r.Lock()
	defer r.Unlock()

	delete(r.DB, id)

	return nil
}
