# 📚 Library Management API

A simple RESTful API to manage library books and borrowing systems.

## 🔌 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **GET** | `/books` | Get all books |
| **GET** | `/books/{book_id}` | Get details of a specific book |
| **POST** | `/books` | Add a new book |
| **PUT** | `/books/{book_id}` | Update book details (Title, Author, Year) |
| **PATCH** | `/books/{book_id}/borrow` | Borrow a book (Business Logic) |
| **PATCH** | `/books/{book_id}/return` | Return a borrowed book |
| **DELETE** | `/books/{book_id}` | Remove a book from the system |

## 🛠️ How to Run

1.  **Prerequisites:**
    * Ensure PostgreSQL is running and create a database named `library_db`.
    * Update the database password in `main.go` inside the `ConnectDB` function.

2.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Start Server:**
    ```bash
    go run main.go
    ```
    The server will start at `http://localhost:8080`.

## 📦 Data Sample (JSON)

**POST /books**
```json
{
  "title": "The Go Programming Language",
  "author": "Alan A. A. Donovan",
  "year": 2015
}