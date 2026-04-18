package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	r.Route("/books", func(r chi.Router) {
		r.Get("/", GetAllBooks)
		r.Post("/", CreateBook)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", GetBook)       // Get books/1
			r.Put("/", UpdateBook)    // Put books/1
			r.Delete("/", DeleteBook) // Delete books/1
		})

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

// Get books/1
func GetBook(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	mu.RLock()
	book, exists := books[id]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(book)

}

// PUT books/1
func UpdateBook(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	var updated Book

	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.RLock()
	if _, exists := books[id]; !exists {
		mu.RUnlock()
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	updated.ID = id
	books[id] = updated
	mu.RUnlock()

	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(updated)
}

// Delete books/1
func DeleteBook(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	mu.RLock()
	delete(books, id)
	mu.RUnlock()

	w.WriteHeader(http.StatusNoContent)
}
