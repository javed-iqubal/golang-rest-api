package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Book struct {
	ID     int    `json.id`
	Title  string `json.title`
	Author string `json.author`
	Year   int    `json.year`
}

var (
	books  = make(map[int]Book)
	nextID = 1
	mu     sync.RWMutex
)

func main() {
	r := chi.NewRouter()

	// route
	r.Route("/book", func(r chi.Router) {
		r.Get("/", GetAllBooks)
		r.Post("/", CreateBook)

	})
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)

}

// ==================== Handlers ====================

// GET /books
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	bookList := make([]Book, 0, len(books))
	for _, book := range books {
		bookList = append(bookList, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookList)
}

// POST /books
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	book.ID = nextID
	nextID++
	books[book.ID] = book
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}
