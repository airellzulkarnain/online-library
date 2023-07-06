# Documentation
## API Endpoints
**1. Create a book**
- URL: POST /books
- Request Body: JSON object representing the book to be created
  - Fields:
    - `title` (string, required): The title of the book.
    - `published_year` (integer, required): The year the book was published.
    - `isbn` (string, required): The ISBN (International Standard Book Number) of the book.
- Response:
  - Status Code: 201 (Created) if successful
  - Response Body: JSON object representing the created book

**2. Get all books**
- URL: GET /books
- Response:
  - Status Code: 200 (OK) if successful
  - Response Body: JSON array containing objects representing all the books
    - Each book object contains the following fields:
      - `id` (unsigned integer): The ID of the book.
      - `title` (string): The title of the book.
      - `published_year` (integer): The year the book was published.
      - `isbn` (string): The ISBN (International Standard Book Number) of the book.

**3. Get a specific book**
- URL: GET /books/:id
- URL Parameters:
  - `id` (unsigned integer): The ID of the book to retrieve.
- Response:
  - Status Code: 200 (OK) if successful
  - Response Body: JSON object representing the retrieved book
    - Fields:
      - `id` (unsigned integer): The ID of the book.
      - `title` (string): The title of the book.
      - `published_year` (integer): The year the book was published.
      - `isbn` (string): The ISBN (International Standard Book Number) of the book.

**4. Update a book**
- URL: PUT /books/:id
- URL Parameters:
  - `id` (unsigned integer): The ID of the book to update.
- Request Body: JSON object representing the updated book data
  - Fields:
    - `title` (string, required): The updated title of the book.
    - `published_year` (integer, required): The updated year the book was published.
    - `isbn` (string, required): The updated ISBN (International Standard Book Number) of the book.
- Response:
  - Status Code: 200 (OK) if successful
  - Response Body: JSON object representing the updated book

**5. Delete a book**
- URL: DELETE /books/:id
- URL Parameters:
  - `id` (unsigned integer): The ID of the book to delete.
- Response:
  - Status Code: 204 (No Content) if successful

**6. Create an author**
- URL: POST /authors
- Request Body: JSON object representing the author to be created
  - Fields:
    - `name` (string, required): The name of the author.
    - `country` (string, required): The country of the author.
- Response:
  - Status Code: 201 (Created) if successful
  - Response Body: JSON object representing the created author

**7. Get all authors**
- URL: GET /authors
- Response:
  - Status Code: 200 (OK) if successful
  - Response Body: JSON array containing objects representing all the authors
    - Each author object contains the following fields:
      - `id` (unsigned integer): The ID of the author.
      - `name` (string): The name of the author.
      - `country` (string): The country of the author.

**8. Get a specific author**
- URL: GET /authors/:id
- URL Parameters:
  - `id` (unsigned integer): The ID of the author to retrieve.
- Response:
  - Status Code: 200 (OK) if successful
  - Response Body: JSON object representing the retrieved author
    - Fields:
      - `id` (unsigned integer): The ID of the author.
      - `name` (string): The name of the author.
      - `country` (string): The country of the author.

**9. Update an author**
- URL: PUT /authors/:id
- URL Parameters:
  - `id` (unsigned integer): The ID of the author to update.
- Request Body: JSON object representing the updated author data
  - Fields:
    - `name` (string, required): The updated name of the author.
    - `country` (string, required): The updated country of the author.
- Response:
  - Status Code: 200 (OK) if successful
  - Response Body: JSON object representing the updated author

**10. Delete an author**
- URL: DELETE /authors/:id
- URL Parameters:
  - `id` (unsigned integer): The ID of the author to delete.
- Response:
  - Status Code: 204 (No Content) if successful

**11. Read all books for a specific author**
- URL: GET /authors/:id/books
- URL Parameters:
  - `id` (unsigned integer): The ID of the author to retrieve.
- Response:
  - Status Code: 200 (OK) if successful
  - Response Body: JSON array containing objects representing all the books
    - Each book object contains the following fields:
      - `id` (unsigned integer): The ID of the book.
      - `title` (string): The title of the book.
      - `published_year` (integer): The year the book was published.
      - `isbn` (string): The ISBN (International Standard Book Number) of the book.

**12. Read all authors for a specific book**
- URL: GET /books/:id/authors
- URL Parameters:
  - `id` (unsigned integer): The ID of the book to retrieve.
- Response:
  - Status Code: 200 (OK) if successful
  - Response Body: JSON array containing objects representing all the authors
    - Each author object contains the following fields:
      - `id` (unsigned integer): The ID of the author.
      - `name` (string): The name of the author.
      - `country` (string): The country of the author.

**13. link book to author**
- URL: POST /books/:book_id/authors/:author_id
- URL Parameters:
  - `book_id` (unsigned integer): The ID of the book to update.
  - `author_id` (unsigned integer): The ID of the author to link the book to.
- Request Body: nil
- Response:
  - Status Code: 204 (NO CONTENT) if successful

## Setup & Running Instructions
### Prerequisite
1. You need git installed on your system.
### Setup
1. clone the repository
```bash
$ git clone https://github.com/airellzulkarnain/online-library.git
```
### Running Instructions
1. Go To bin directory.
Linux/Mac: 
```bash
$ cd online-library/bin
```
Windows:
```bash
$ cd online-library\\bin
```
2. Run the online-library executable.
Linux/Mac:
```bash
$ ./online-library
```
Windows:
```bash
$ online-library.exe
```
Now you get the API up and running on http://localhost:8080/
3. Before you can access any operations, you must get jwt token through the `login/` endpoint:
```bash
curl -X POST -d '{"username": "zegen", "password": "zegen"}' http://localhost:8080/login
```
output:
```bash
{"token": "jwt-token"}
```
### Operations
**1. Create a book**
```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{
  "title": "The Great Gatsby",
  "published_year": 1925,
  "isbn": 123456789011X
}' http://localhost:8080/api/books
```

**2. Get all books**
```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/books
```

**3. Get a specific book**
```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/books/1
```

**4. Update a book**
```bash
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{
  "title": "The Great Gatsby",
  "published_year": 1925,
  "isbn": "123456789012X"
}' http://localhost:8080/api/books/1
```

**5. Delete a book**
```bash
curl -X DELETE -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/books/1
```

**6. Create an author**
```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{
  "name": "F. Scott Fitzgerald",
  "country": "United States"
}' http://localhost:8080/api/authors
```

**7. Get all authors**
```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/authors
```

**8. Get a specific author**
```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/authors/1
```

**9. Update an author**
```bash
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{
  "name": "F. Scott Fitzgerald",
  "country": "United States"
}' http://localhost:8080/api/authors/1
```

**10. Delete an author**
```bash
curl -X DELETE -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/authors/1
```

**11. Get all books for a specific author**
```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/authors/1/books
```

**12. Get all authors specific books**
```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/books/1/authors
```
**13. Link book to author**
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/books/1/authors/1/
```

Make sure to replace `$TOKEN` with the actual JWT token value you get from login endpoint earlier.

## Testing Procedures
