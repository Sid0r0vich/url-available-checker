package handlers

import (
	"net/http"

	"github.com/Sid0r0vich/url-available-checker/internal/repository"
)

type API struct {
	Repo repository.Repository
}

func NewAPI() *API {
	return &API{Repo: repository.NewInMemoryRepo()}
}

func NewMux(api *API) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/links", api.GetLinksHanlder)
	r.HandleFunc("/list", api.MakePDFHandler)

	return r
}
