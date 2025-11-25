package main

import (
	"fmt"
	"net/http"

	handlers "github.com/Sid0r0vich/url-available-checker/internal/handlers"
)

func main() {
	r := http.NewServeMux()
	api := handlers.NewAPI()

	r.HandleFunc("/links", api.GetLinksHanlder)
	r.HandleFunc("/list", api.MakePDFHandler)

	if err := http.ListenAndServe(":8082", r); err != nil {
		fmt.Printf("Server error: %v", err)
	}
}
