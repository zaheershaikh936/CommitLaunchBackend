package config

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ServerConfig(route *chi.Mux) error {
	fmt.Println("server is started on http://localhost:8080")
	return http.ListenAndServe(":8080", route)
}
