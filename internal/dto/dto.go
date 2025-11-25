package dto

type Link struct {
	URL          string
	Availability bool
}

type LinksRequest struct {
	Links []string `json:"links"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
