package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Sid0r0vich/url-available-checker/internal/dto"
)

type check struct {
	Link dto.Link
	Ind  int
}

func checkURL(url string, client *http.Client, result chan<- *check, ind int) {
	resp, err := client.Get("https://" + url)
	if err == nil {
		resp.Body.Close()
	}
	result <- &check{Link: dto.Link{URL: url, Availability: err == nil}, Ind: ind}
}

func (api *API) GetLinksHanlder(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	var req dto.LinksRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(dto.ErrorResponse{Error: "error parsing json"})
		return
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	linkAvailability := make([]dto.Link, len(req.Links))
	respLinks := make(map[string]string, len(req.Links))
	result := make(chan *check)
	var wg sync.WaitGroup
	for ind, url := range req.Links {
		wg.Add(1)
		go func() {
			checkURL(url, client, result, ind)
			wg.Done()
		}()
	}

	var wgConsumer sync.WaitGroup
	wgConsumer.Add(1)
	go func() {
		for {
			checkedURL, ok := <-result
			if !ok {
				break
			}

			linkAvailability[checkedURL.Ind] = checkedURL.Link
			if checkedURL.Link.Availability {
				respLinks[checkedURL.Link.URL] = "available"
			} else {
				respLinks[checkedURL.Link.URL] = "not available"
			}
		}

		wgConsumer.Done()
	}()

	wg.Wait()
	close(result)
	wgConsumer.Wait()

	id := api.Repo.AddLinks(linkAvailability)
	w.WriteHeader(http.StatusOK)
	enc.Encode(dto.LinkResponse{Links: respLinks, Num: id})
}

func (api *API) MakePDFHandler(w http.ResponseWriter, r *http.Request) {

}
