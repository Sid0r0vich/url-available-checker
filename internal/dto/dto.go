package dto

type Link struct {
	URL          string
	Availability bool
}

type LinksRequest struct {
	Links []string `json:"links"`
}

type MakePDFRequest struct {
	LinksList []int `json:"links_list"`
}

type LinkResponse struct {
	Links map[string]string `json:"links"`
	Num   int               `json:"links_num"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{Error: msg}
}
