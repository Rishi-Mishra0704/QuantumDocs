package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

var books = []Book{
	{ID: "1", Title: "The Go Programming Language", Author: "Alan A. A. Donovan", Price: 34.99},
	{ID: "2", Title: "Clean Code", Author: "Robert C. Martin", Price: 29.99},
}

// @Method GET
// @Path /books/{id}
// @Description Get all books
// @Response 200 - Returns an array of Book objects.
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// @Method GET
// @Path /books/{id}
// @Description Get a book by ID
// @Parameter id - string - The ID of the book to retrieve
// @Response 200 - Returns a single Book object
// @Response 404 - Book not found
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Path[len("/books/"):]
	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.NotFound(w, r)
}

// @Method POST
// @Path /books
// @Description Add a new book
// @Parameter Book - JSON object - The book to add
// @Response 201 - Returns the created Book object
// @Response 400 - Invalid request body
func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	books = append(books, book)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func main() {
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBooks(w, r)
		case http.MethodPost:
			addBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/", getBook)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
