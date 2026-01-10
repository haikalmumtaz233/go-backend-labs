package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title      string `json:"title" gorm:"type:varchar(100);not null"`
	Author     string `json:"author" gorm:"type:varchar(100);not null"`
	Year       int    `json:"year"`
	IsBorrowed bool   `json:"is_borrowed" gorm:"default:false"`
}

var db *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=your_pgadmin_password dbname=library_db port=5432 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Library Database", err.Error())
	}

	db.AutoMigrate(&Book{})
	fmt.Println("Successfully Connected to the Library Database!")
}

// HANDLERS

func listBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book

	if result := db.Find(&books); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBookbyId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("book_id")
	var book Book

	if err := db.First(&book, id).Error; err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook Book

	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if newBook.Title == "" || newBook.Author == "" {
		http.Error(w, "Title or Author cannot be blank", http.StatusBadRequest)
		return
	}

	newBook.IsBorrowed = false

	if result := db.Create(&newBook); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "Success",
		"message": "Book succesfully added",
		"data":    newBook,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("book_id")

	var book Book
	var updatedBook Book

	if err := db.First(&book, id).Error; err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	db.Model(&book).Updates(updatedBook)

	response := map[string]interface{}{
		"status":  "Success",
		"message": "Book has been updated",
		"data":    book,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func borrowBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("book_id")
	var book Book

	if err := db.First(&book, id).Error; err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	if book.IsBorrowed {
		http.Error(w, "This book is borrowed!", http.StatusBadRequest)
		return
	}

	book.IsBorrowed = true
	db.Save(&book)

	response := map[string]interface{}{
		"status":  "Success",
		"message": "Successfully borrow book",
		"data":    book,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func returnBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("book_id")
	var book Book

	if err := db.First(&book, id).Error; err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	book.IsBorrowed = false
	db.Save(&book)

	response := map[string]interface{}{
		"status":  "Success",
		"message": "Successfully return book",
		"data":    book,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("book_id")
	var book Book

	if err := db.First(&book, id).Error; err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	db.Delete(&book)
	response := map[string]interface{}{
		"status":  "Success",
		"message": "Successfully delete book",
		"data":    book,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	ConnectDB()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", listBooks)
	mux.HandleFunc("GET /books/{book_id}", getBookbyId)
	mux.HandleFunc("POST /books", createBook)

	mux.HandleFunc("PUT /books/{book_id}", updateBook)
	mux.HandleFunc("PATCH /books/{book_id}/borrow", borrowBook)
	mux.HandleFunc("PATCH /books/{book_id}/return", returnBook)

	mux.HandleFunc("DELETE /books/{book_id}", deleteBook)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
