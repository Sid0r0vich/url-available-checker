package handlers

import "github.com/Sid0r0vich/url-available-checker/internal/repository"

type API struct {
	Repo repository.Repository
}

func NewAPI() *API {
	return &API{Repo: repository.NewInMemoryRepo()}
}
