package main

import (
	"net/http"

	handlers "github.com/Sid0r0vich/url-available-checker/internal/handlers"
)

func main() {
	r := http.NewServeMux()
	api := handlers.NewAPI()

	r.HandleFunc("/linkss", api.GetLinksHanlder)
	r.HandleFunc("/list", api.MakePDFHandler)
}
