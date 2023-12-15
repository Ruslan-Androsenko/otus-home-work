package internalhttp

import (
	"net/http"
)

// Обработчик для главной страницы.
func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	_, err := w.Write([]byte("Welcome to the Home Page!"))
	if err != nil {
		logg.Errorf("Failed to display home page. Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Обработчик для приветственной страницы.
func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Welcome to the Hello Page!"))
	if err != nil {
		logg.Errorf("Failed to display hello page. Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
