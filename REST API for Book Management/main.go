package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear string `json:"publication_year"`
}

var (
	books = make(map[int]Book)
	mutex = &sync.Mutex{}
)

func main() {

	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBook(w, r)
		case http.MethodPost:
			createBook(w, r)
		case http.MethodPut:
			updateBook(w, r)
		case http.MethodDelete:
			deleteBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080",nil)
}

/*  Functions of CRUD  */

// get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(books)
}

// get one book by id
func getBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID ", http.StatusBadRequest)
		return
	}

	book, exists := books[id]
	if !exists {
		http.Error(w, "Book not found ", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	if book.ID == 0 {
		http.Error(w, "Book id is required", http.StatusBadRequest)
		return
	}

	books[book.ID] = book

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateBook Book
	if err := json.NewDecoder(r.Body).Decode(&updateBook); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	if updateBook.ID == 0 {
		http.Error(w, "Book id do not match", http.StatusBadRequest)
		return
	}

	books[id] = updateBook

	w.Header().Set("Content-Type", "application/josn")
	json.NewEncoder(w).Encode(updateBook)
}

// delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalide ID", http.StatusBadRequest)
		return
	}

	if _, exists := books[id]; !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	delete(books, id)

	w.WriteHeader(http.StatusNoContent)
}
