package memory

import (
	"context"
	"github.com/koind/proxy/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

const requestID = "23534lkj3lk4j55"

var requestRepository repository.RequestRepositoryInterface

func init() {
	requestRepository = NewRequestRepository()
}

func before() {
	req := repository.Row{
		ID:       requestID,
		Request:  repository.Request{},
		Response: repository.Response{},
	}

	requestRepository.Create(context.Background(), req)
}

func after() {
	requestRepository.Delete(context.Background(), requestID)
}

func TestRequestRepository_Create(t *testing.T) {
	newRequest := repository.Row{
		ID:       requestID,
		Request:  repository.Request{},
		Response: repository.Response{},
	}

	_, err := requestRepository.Create(context.Background(), newRequest)
	assert.Nil(t, err, "не должно быть ошибки при создании")

	req, _ := requestRepository.GetOneByID(context.Background(), requestID)
	assert.EqualValues(t, &newRequest, req)

	after()
}

func TestRequestRepository_GetOneById(t *testing.T) {
	before()

	_, err := requestRepository.GetOneByID(context.Background(), requestID)
	assert.Nil(t, err, "не должно быть ошибки при получении записи")

	after()

	_, err = requestRepository.GetOneByID(context.Background(), "test-random-id")
	if assert.NotNil(t, err) {
		assert.Equal(t, err, repository.ErrNotFount, "ошибки должны совподать")
	}
}

func TestRequestRepository_FindAll(t *testing.T) {
	before()

	list, err := requestRepository.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	after()

	list, _ = requestRepository.GetAll(context.Background())
	assert.Len(t, list, 0)
}

func TestRequestRepository_Delete(t *testing.T) {
	before()

	err := requestRepository.Delete(context.Background(), requestID)
	assert.Nil(t, err, "не должно быть ошибки при удалении")

	_, err = requestRepository.GetOneByID(context.Background(), requestID)
	if assert.NotNil(t, err) {
		assert.Equal(t, err, repository.ErrNotFount, "ошибки должны совподать")
	}
}
