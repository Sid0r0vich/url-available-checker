package repository

import (
	"errors"
	"sync"

	"github.com/Sid0r0vich/url-available-checker/internal/dto"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	AddLinks([]dto.Link)
	GetLinks([]int) ([]dto.Link, error)
}

type InMemoryRepo struct {
	sync.RWMutex
	IdToLink map[int]dto.Link
	URLToId  map[string]int
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{IdToLink: make(map[int]dto.Link), URLToId: make(map[string]int)}
}

func (repo *InMemoryRepo) AddLinks(links []dto.Link) {
	repo.Lock()
	defer repo.Unlock()

	for _, link := range links {
		var newId int
		if id, ok := repo.URLToId[link.URL]; ok {
			newId = id
		} else {
			newId = len(repo.IdToLink)
		}
		repo.IdToLink[newId] = link
	}
}

func (repo *InMemoryRepo) GetLinks(ids []int) ([]dto.Link, error) {
	repo.RLock()
	defer repo.RUnlock()

	links := make([]dto.Link, len(ids))
	for ind, id := range ids {
		link, ok := repo.IdToLink[id]
		if !ok {
			return nil, ErrNotFound
		}

		links[ind] = link
	}

	return links, nil
}
