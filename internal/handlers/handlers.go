package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Sid0r0vich/url-available-checker/internal/dto"
	"github.com/Sid0r0vich/url-available-checker/internal/pdf"
	"github.com/Sid0r0vich/url-available-checker/internal/repository"
)

type check struct {
	Link dto.Link
	Ind  int
}

func checkURL(url string, client *http.Client, result chan<- *check, ind int) {
	req, err := http.NewRequest("GET", "https://"+url, nil)
	if err != nil {
		fmt.Printf("request error: %v", err)
		result <- &check{Link: dto.Link{URL: url, Availability: false}, Ind: ind}
		return
	}

	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
	result <- &check{Link: dto.Link{URL: url, Availability: err == nil}, Ind: ind}
}

func linksToPrettyLinks(links []dto.Link) map[string]string {
	res := make(map[string]string, len(links))
	for _, link := range links {
		if link.Availability {
			res[link.URL] = "available"
		} else {
			res[link.URL] = "not available"
		}
	}

	return res
}

func (api *API) GetLinksHanlder(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	select {
	case <-api.Context.Done():
		w.WriteHeader(http.StatusServiceUnavailable)
		enc.Encode(dto.NewErrorResponse("server unavailable"))
		return
	default:
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		enc.Encode(dto.NewErrorResponse("method not allowed"))
		return
	}

	var req dto.LinksRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(dto.NewErrorResponse("error parsing json"))
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	result := make(chan *check)
	var wg sync.WaitGroup
	for ind, url := range req.Links {
		wg.Add(1)
		go func(u string, i int) {
			checkURL(u, client, result, i)
			wg.Done()
		}(url, ind)
	}

	linkAvailability := make([]dto.Link, len(req.Links))
	var wgConsumer sync.WaitGroup
	wgConsumer.Add(1)
	go func() {
		for {
			checkedURL, ok := <-result
			if !ok {
				break
			}

			linkAvailability[checkedURL.Ind] = checkedURL.Link
		}

		wgConsumer.Done()
	}()

	wg.Wait()
	close(result)
	wgConsumer.Wait()

	id := api.Repo.AddLinks(linkAvailability)
	w.WriteHeader(http.StatusOK)
	enc.Encode(dto.LinkResponse{Links: linksToPrettyLinks(linkAvailability), Num: id})
}

func (api *API) MakePDFHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		enc.Encode(dto.NewErrorResponse("method not allowed"))
		return
	}

	var req dto.MakePDFRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(dto.NewErrorResponse("error parsing json"))
		return
	}

	links, err := api.Repo.GetLinks(req.LinksList)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			enc.Encode(dto.NewErrorResponse("links_num not found"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(dto.NewErrorResponse("unexpected error"))
		}

		return
	}

	buf, err := pdf.GeneratePDFFromLinks(links)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(dto.NewErrorResponse("unexpected error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="links.pdf"`)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("error writing body: %v\n", err)
	}
}
