package dto

import "time"

type Link struct {
	URL          string
	Availability bool
	Time         time.Time
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
