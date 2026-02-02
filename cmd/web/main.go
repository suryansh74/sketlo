package main

import (
	"net/http"
)

func main() {
	handler := SetupRoutes()
	http.ListenAndServe(":8000", handler)
}
