package dto

type Link struct {
	URL          string
	Availability bool
}

type LinksRequest struct {
	Links []string `json:"links"`
}

type LinkResponse struct {
	Links map[string]string `json:"links"`
	Num   int               `json:"links_num"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
