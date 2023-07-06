package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	PublishedYear int    `json:"published_year"`
	ISBN          string `json:"isbn"`
}

type Author struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

var (
	db  *sql.DB
	err error
)

// Create the books and authors tables
func createTables() {
	booksTableSQL := `
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			published_year INTEGER NOT NULL,
			isbn TEXT NOT NULL
		);`
	_, err = db.Exec(booksTableSQL)
	if err != nil {
		log.Fatal("Failed to create books table:", err)
	}

	authorsTableSQL := `
		CREATE TABLE IF NOT EXISTS authors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			country TEXT NOT NULL
		);`
	_, err = db.Exec(authorsTableSQL)
	if err != nil {
		log.Fatal("Failed to create authors table:", err)
	}

	booksAuthorsTableSQL := `
		CREATE TABLE IF NOT EXISTS books_authors (
			book_id INTEGER NOT NULL,
			author_id INTEGER NOT NULL,
			FOREIGN KEY(book_id) REFERENCES books(id),
			FOREIGN KEY(author_id) REFERENCES authors(id)
		);`
	_, err = db.Exec(booksAuthorsTableSQL)
	if err != nil {
		log.Fatal("Failed to create books_authors table:", err)
	}
}

// Auth middleware
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.Split(tokenString, "Bearer ")[1]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("JtQmEYnaYDj476+w+NmsXwWS8sBcftCgVwuhupDK+YW9ohM7W/mi+BM7n3uxaKL9Z1p5OQ4Ory634Yz7"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}

// Handlers

func createBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if book.Title == "" || book.PublishedYear == 0 || book.ISBN == "" || len(book.ISBN) != 13 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Create the book
	stmt, err := db.Prepare("INSERT INTO books (title, published_year, isbn) VALUES (?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}
	defer stmt.Close()

	var r sql.Result
	r, err = stmt.Exec(book.Title, book.PublishedYear, book.ISBN)
	id, _ := r.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	book.ID = uint(id)
	c.JSON(http.StatusCreated, book)
}

func getBooks(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, published_year, isbn FROM books")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.PublishedYear, &book.ISBN); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
			return
		}
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id := c.Param("id")

	var book Book
	err := db.QueryRow("SELECT id, title, published_year, isbn FROM books WHERE id = ?", id).
		Scan(&book.ID, &book.Title, &book.PublishedYear, &book.ISBN)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func updateBook(c *gin.Context) {
	id := c.Param("id")

	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if book.Title == "" || book.PublishedYear == 0 || book.ISBN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Update the book
	stmt, err := db.Prepare("UPDATE books SET title = ?, published_year = ?, isbn = ? WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(book.Title, book.PublishedYear, book.ISBN, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")

	stmt, err := db.Prepare("DELETE FROM books WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func createAuthor(c *gin.Context) {
	var author Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if author.Name == "" || author.Country == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Create the author
	stmt, err := db.Prepare("INSERT INTO authors (name, country) VALUES (?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}
	defer stmt.Close()
	var r sql.Result
	r, err = stmt.Exec(author.Name, author.Country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}
	id, _ := r.LastInsertId()
	author.ID = uint(id)
	c.JSON(http.StatusCreated, author)
}

func getAuthors(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, country FROM authors")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve authors"})
		return
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var author Author
		if err := rows.Scan(&author.ID, &author.Name, &author.Country); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve authors"})
			return
		}
		authors = append(authors, author)
	}

	c.JSON(http.StatusOK, authors)
}

func getAuthor(c *gin.Context) {
	id := c.Param("id")

	var author Author
	err := db.QueryRow("SELECT id, name, country FROM authors WHERE id = ?", id).
		Scan(&author.ID, &author.Name, &author.Country)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve author"})
		return
	}

	c.JSON(http.StatusOK, author)
}

func updateAuthor(c *gin.Context) {
	id := c.Param("id")

	var author Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if author.Name == "" || author.Country == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Update the author
	stmt, err := db.Prepare("UPDATE authors SET name = ?, country = ? WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(author.Name, author.Country, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	c.JSON(http.StatusOK, author)
}

func deleteAuthor(c *gin.Context) {
	id := c.Param("id")

	stmt, err := db.Prepare("DELETE FROM authors WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func getBooksByAuthor(c *gin.Context) {
	authorID := c.Param("id")

	rows, err := db.Query(`SELECT b.id, b.title, b.published_year, b.isbn FROM books AS b
							INNER JOIN books_authors AS ba ON b.id = ba.book_id
							INNER JOIN authors AS a ON a.id = ba.author_id
							WHERE a.id = ?`, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books by author"})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.PublishedYear, &book.ISBN); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books by author"})
			return
		}
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

func getAuthorsByBook(c *gin.Context) {
	bookID := c.Param("id")

	rows, err := db.Query(`SELECT a.id, a.name, a.country FROM authors AS a
							INNER JOIN books_authors AS ba ON a.id = ba.author_id
							INNER JOIN books AS b ON b.id = ba.book_id
							WHERE b.id = ?`, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve authors by book"})
		return
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var author Author
		if err := rows.Scan(&author.ID, &author.Name, &author.Country); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve authors by book"})
			return
		}
		authors = append(authors, author)
	}

	c.JSON(http.StatusOK, authors)
}

func linkBookToAuthor(c *gin.Context) {
	BookID := c.Param("book_id")
	AuthorID := c.Param("author_id")

	stmt, err := db.Prepare("INSERT INTO books_authors (book_id, author_id) VALUES (?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link book to author"})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(BookID, AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link book to author"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user ()
	if user.Username != "zegen" || user.Password != "zegen" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	tokenString, err := token.SignedString([]byte("JtQmEYnaYDj476+w+NmsXwWS8sBcftCgVwuhupDK+YW9ohM7W/mi+BM7n3uxaKL9Z1p5OQ4Ory634Yz7"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func main() {
	// Initialize database
	db, _ = sql.Open("sqlite3", "./library.db?_foreign_keys=on")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer db.Close()

	// Create books table
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS books (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						title TEXT NOT NULL,
						published_year INTEGER NOT NULL,
						isbn TEXT NOT NULL UNIQUE
					)`)

	// Create authors table
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS authors (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT NOT NULL,
						country TEXT NOT NULL, 
						CONSTRAINT UC_name_country UNIQUE (name, country)
					)`)

	// Create books_authors table
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS books_authors (
						book_id INTEGER,
						author_id INTEGER,
						FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE,
						FOREIGN KEY (author_id) REFERENCES authors (id) ON DELETE CASCADE,
						PRIMARY KEY (book_id, author_id)
					)`)

	r := gin.Default()

	// Public routes
	r.POST("/login", login)

	// Protected routes
	api := r.Group("/api")
	api.Use(authMiddleware())

	{
		api.GET("/books", getBooks)
		api.POST("/books", createBook)
		api.GET("/books/:id", getBook)
		api.PUT("/books/:id", updateBook)
		api.DELETE("/books/:id", deleteBook)

		api.GET("/authors", getAuthors)
		api.POST("/authors", createAuthor)
		api.GET("/authors/:id", getAuthor)
		api.PUT("/authors/:id", updateAuthor)
		api.DELETE("/authors/:id", deleteAuthor)

		api.POST("/books/:book_id/authors/:author_id", linkBookToAuthor)
		api.GET("/authors/:id/books", getBooksByAuthor)
		api.GET("/books/:id/authors", getAuthorsByBook)
	}

	r.Run(":8080")
}
