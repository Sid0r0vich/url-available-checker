package handlers

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sid0r0vich/url-available-checker/internal/repository"
)

type API struct {
	Repo    repository.Repository
	Context context.Context
	Cancel  context.CancelFunc
}

func NewAPI() *API {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	return &API{
		Repo:    repository.NewInMemoryRepo(),
		Context: ctx,
		Cancel:  cancel,
	}
}
