package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var router *gin.Engine

func teardown() {
	// Close the database connection
	db.Close()
}

func setupRoutes() {
	router.POST("/books", createBook)
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBook)
	router.PUT("/books/:id", updateBook)
	router.DELETE("/books/:id", deleteBook)
	router.POST("/authors", createAuthor)
	router.GET("/authors", getAuthors)
	router.GET("/authors/:id", getAuthor)
	router.PUT("/authors/:id", updateAuthor)
	router.DELETE("/authors/:id", deleteAuthor)
}

func TestCreateBook(t *testing.T) {
	// Create a test HTTP request
	requestBody := []byte(`{"title": "Book 1", "published_year": 2022, "isbn": "123456789011x"}`)
	request, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	// Create a test HTTP response recorder
	recorder := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status 201, but got %d", recorder.Code)
	}

	// Check the response body
	expectedResponseBody := `{"id":1,"title":"Book 1","published_year":2022,"isbn":"123456789011x"}`
	if recorder.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponseBody, recorder.Body.String())
	}
}

func TestGetBooks(t *testing.T) {
	// Insert a book into the database for testing
	insertBook("Book 1", 2022, "123456789011x")

	// Create a test HTTP request
	request, _ := http.NewRequest("GET", "/books", nil)

	// Create a test HTTP response recorder
	recorder := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", recorder.Code)
	}

	// Check the response body
	expectedResponseBody := `[{"id":1,"title":"Book 1","published_year":2022,"isbn":"123456789011x"}]`
	if recorder.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponseBody, recorder.Body.String())
	}
}

func TestGetBook(t *testing.T) {
	// Insert a book into the database for testing
	insertBook("Book 1", 2022, "123456789011x")

	// Create a test HTTP request
	request, _ := http.NewRequest("GET", "/books/1", nil)

	// Create a test HTTP response recorder
	recorder := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", recorder.Code)
	}

	// Check the response body
	expectedResponseBody := `{"id":1,"title":"Book 1","published_year":2022,"isbn":"123456789011x"}`
	if recorder.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponseBody, recorder.Body.String())
	}
}

// create TestUpdateBook
func TestUpdateBook(t *testing.T) {
	// Insert a book into the database for testing
	insertBook("Book 1", 2022, "123456789011x")

	// Create a test HTTP request
	request, _ := http.NewRequest("PUT", "/books/1", bytes.NewBuffer([]byte(`{"title": "Book 2", "published_year": 2023, "isbn": "123456789012x"}`)))
	request.Header.Set("Content-Type", "application/json")

	// Create a test HTTP response recorder
	recorder := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", recorder.Code)
	}

}

func TestDeleteBook(t *testing.T) {
	// Insert a book into the database for testing
	insertBook("Book 1", 2022, "123456789011x")
	r, _ := http.NewRequest("DELETE", "/books/1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, r)
	if recorder.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, but got %d", recorder.Code)
	}
}

// Helper function to insert a book into the database
func insertBook(title string, publishedYear int, isbn string) {
	stmt, _ := db.Prepare("DELETE FROM books;INSERT INTO books (title, published_year, isbn) VALUES (?, ?, ?)")
	stmt.Exec(title, publishedYear, isbn)
}

// create insertAuthor with name and country param
func insertAuthor(name string, country string) {
	stmt, _ := db.Prepare("INSERT INTO authors (name, country) VALUES (?, ?)")
	stmt.Exec(name, country)
}

// create TestGetAuthors same as TestGetBooks
func TestGetAuthors(t *testing.T) {
	// Insert a book into the database for testing
	insertAuthor("Author 1", "USA")

	// Create a test HTTP request
	request, _ := http.NewRequest("GET", "/authors", nil)

	// Create a test HTTP response recorder
	recorder := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", recorder.Code)
	}

	// Check the response body
	expectedResponseBody := `[{"id":1,"name":"Author 1","country":"USA"}]`
	if recorder.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponseBody, recorder.Body.String())
	}
}

// create TestGetAuthor
func TestGetAuthor(t *testing.T) {
	// Insert a book into the database for testing
	insertAuthor("Author 1", "USA")

	// Create a test HTTP request
	request, _ := http.NewRequest("GET", "/authors/1", nil)

	// Create a test HTTP response recorder
	recorder := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", recorder.Code)
	}

	// Check the response body
	expectedResponseBody := `{"id":1,"name":"Author 1","country":"USA"}`
	if recorder.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponseBody, recorder.Body.String())
	}

}

// Create TestUpdateAuthor
func TestUpdateAuthor(t *testing.T) {
	// Insert a book into the database for testing
	insertAuthor("Author 1", "USA")

	// Create a test HTTP request
	request, _ := http.NewRequest("PUT", "/authors/1", bytes.NewBuffer([]byte(`{"name": "Author 2", "country": "France"}`)))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", recorder.Code)
	}
}

// create TestDeleteAuthor
func TestDeleteAuthor(t *testing.T) {
	// Insert a book into the database for testing
	insertAuthor("Author 1", "USA")

	// Create a test HTTP request
	request, _ := http.NewRequest("DELETE", "/authors/1", nil)

	// Create a test HTTP response recorder
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, but got %d", recorder.Code)
	}
}

func TestLinkBookToAuthor(t *testing.T) {
	// insertbook
	insertBook("Book 1", 2022, "123456789011x")
	// insertauthor
	insertAuthor("Author 1", "USA")

	// Create a test HTTP request
	request, _ := http.NewRequest("POST", "/books/1/authors/1", nil)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, but got %d", recorder.Code)
	}
}

func TestIntegration(t *testing.T) {
	// Run integration tests
	// Note: This is just an example to demonstrate integration testing
	// In a real-world scenario, you may want to set up and tear down the test environment differently
	setup()
	defer teardown()

	// Create a book
	requestBodyBook := []byte(`{"title": "Book 1", "published_year": 2022, "isbn": "123456789011x"}`)
	requestBook, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(requestBodyBook))
	requestBook.Header.Set("Content-Type", "application/json")
	recorderBook := httptest.NewRecorder()
	router.ServeHTTP(recorderBook, requestBook)

	// Verify the response
	if recorderBook.Code != http.StatusCreated {
		t.Errorf("Expected status 201, but got %d", recorderBook.Code)
	}

	// Get the created book
	requestBook, _ = http.NewRequest("GET", "/books/1", nil)
	recorderBook = httptest.NewRecorder()
	router.ServeHTTP(recorderBook, requestBook)

	// Verify the response
	expectedResponseBody := `{"id":1,"title":"Book 1","published_year":2022,"isbn":"123456789011x"}`
	if recorderBook.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponseBody, recorderBook.Body.String())
	}

	// create an author
	requestBodyAuthor := []byte(`{"name": "Author 1", "country": "USA"}`)
	requestAuthor, _ := http.NewRequest("POST", "/authors", bytes.NewBuffer(requestBodyAuthor))
	requestAuthor.Header.Set("Content-Type", "application/json")
	recorderAuthor := httptest.NewRecorder()
	router.ServeHTTP(recorderAuthor, requestAuthor)

	// Verify the response
	if recorderAuthor.Code != http.StatusCreated {
		t.Errorf("Expected status 201, but got %d", recorderAuthor.Code)
	}

	// Get the created author
	requestAuthor, _ = http.NewRequest("GET", "/authors/1", nil)
	recorderAuthor = httptest.NewRecorder()
	router.ServeHTTP(recorderAuthor, requestAuthor)

	// Verify the response
	expectedResponseBody = `{"id":1,"name":"Author 1","country":"USA"}`
	if recorderAuthor.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponseBody, recorderAuthor.Body.String())
	}

	var (
		book   Book
		author Author
	)

	// decode recorderBook.Body.String() into book
	if err := json.NewDecoder(recorderBook.Body).Decode(&book); err != nil {
		t.Fatal(err)
	}

	// same as above for author
	if err := json.NewDecoder(recorderAuthor.Body).Decode(&author); err != nil {
		t.Fatal(err)
	}

	// link book to author
	requestLink, _ := http.NewRequest("POST", "/books/1/authors/1", nil)
	requestLink.Header.Set("Content-Type", "application/json")
	recorderLink := httptest.NewRecorder()
	router.ServeHTTP(recorderLink, requestLink)

	// Verify the response
	if recorderLink.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, but got %d", recorderLink.Code)
	}

	// delete book
	requestDelete, _ := http.NewRequest("DELETE", "/books/1", nil)
	recorderDelete := httptest.NewRecorder()
	router.ServeHTTP(recorderDelete, requestDelete)

	// Verify the response
	if recorderDelete.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, but got %d", recorderDelete.Code)
	}

	// delete author
	requestDelete, _ = http.NewRequest("DELETE", "/authors/1", nil)
	recorderDelete = httptest.NewRecorder()
	router.ServeHTTP(recorderDelete, requestDelete)

	// Verify the response
	if recorderDelete.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, but got %d", recorderDelete.Code)
	}

}

func setup() {
	// Initialize the Gin router
	router = gin.Default()

	// Connect to the in-memory SQLite database for testing
	db, _ = sql.Open("sqlite3", ":memory:")

	// Create the tables for testing
	createTables()

	// Set the router to use the test database
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Set up the API routes for testing
	setupRoutes()
}

func TestMain(m *testing.M) {
	// Set up test environment
	setup()

	// Run tests
	exitCode := m.Run()

	// Tear down test environment
	teardown()

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}
