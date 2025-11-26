package main

import (
	"fmt"
	"net/http"

	handlers "github.com/Sid0r0vich/url-available-checker/internal/handlers"
)

var (
	PORT = "8082"
)

func main() {
	r := http.NewServeMux()
	api := handlers.NewAPI()
	defer api.Cancel()

	r.HandleFunc("/links", api.GetLinksHanlder)
	r.HandleFunc("/list", api.MakePDFHandler)

	fmt.Printf("Server started on :%s\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r); err != nil {
		fmt.Printf("Server error: %v", err)
	}
}
