# 📚 Library Management API

A simple RESTful API to manage library books and borrowing systems.

## 🚀 Key Features

* **CRUD Operations:** Create, Read, Update, Delete books.
* **Security:**
    * **Anti-SQL Injection:** Uses GORM parameterized queries for searching.
    * **Environment Variables:** Sensitive data managed via `.env`.
* **Soft Delete System:**
    * **Soft Delete:** Data is marked as deleted but kept in the DB (`deleted_at`).
    * **Restore:** Ability to bring back soft-deleted data.
    * **Hard Delete:** Permanently remove data from the database.
* **Business Logic:** Borrowing and Returning mechanism.

## 🔌 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **GET** | `/books` | Get all active books |
| **GET** | `/books/{id}` | Get specific book detail |
| **GET** | `/books/search?q=...` | **Search** books by Title or Author |
| **POST** | `/books` | Add a new book |
| **PUT** | `/books/{id}` | Update book details |
| **PATCH** | `/books/{id}/borrow` | Borrow a book |
| **PATCH** | `/books/{id}/return` | Return a borrowed book |
| **DELETE** | `/books/{id}` | **Soft Delete** (Move to Trash) |
| **DELETE** | `/books/{id}/permanent` | **Hard Delete** (Remove Forever) |
| **PATCH** | `/books/{id}/restore` | **Restore** deleted book |

## 🛠️ How to Run

1.  **Prerequisites:**
    * Ensure PostgreSQL is running and create a database named `library_db`.
    * Update the database password in `main.go` inside the `ConnectDB` function.

2.  **Environment Setup:**
    Create a `.env` file in the root directory:
    ```env
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_NAME=library_db
    DB_PORT=5432
    ```

3.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

4.  **Start Server:**
    ```bash
    go run main.go
    ```
    Server runs at `http://localhost:8080`.

## 📦 Data Sample (JSON)

**POST /books**
```json
{
  "title": "The Library API",
  "author": "Haikal Mumtaz",
  "publisher": "Radz",
  "year": 2025
}