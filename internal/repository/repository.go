package repository

import (
	"errors"
	"sync"

	"github.com/Sid0r0vich/url-available-checker/internal/dto"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	AddLinks([]dto.Link) int
	GetLinks([]int) ([]dto.Link, error)
}

type InMemoryRepo struct {
	sync.RWMutex
	Storage map[int][]dto.Link
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{Storage: make(map[int][]dto.Link)}
}

func (repo *InMemoryRepo) AddLinks(links []dto.Link) int {
	repo.Lock()
	defer repo.Unlock()

	id := len(repo.Storage)
	repo.Storage[id] = links

	return id
}

func (repo *InMemoryRepo) GetLinks(ids []int) ([]dto.Link, error) {
	repo.RLock()
	defer repo.RUnlock()

	links := make([]dto.Link, 0)
	for _, id := range ids {
		ls, ok := repo.Storage[id]
		if !ok {
			return nil, ErrNotFound
		}

		links = append(links, ls...)
	}

	return links, nil
}
