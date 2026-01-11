# 📦 Warehouse Management API (WMS Lite)

A simple RESTful API to manage warehouse inventory and stock movements.

## 🚀 Key Features

* **CRUD Operations:** Create, Read, Update, Delete products.
* **Security:**
    * **Anti-SQL Injection:** Uses GORM parameterized queries for product searching.
    * **Environment Variables:** Sensitive data managed via `.env`.
* **Smart Stock Management:**
    * **Stock Adjustment:** Dedicated logic for handling Stock In/Out (prevents data race conditions).
    * **Validation:** Prevents negative stock during outgoing adjustments.
* **Soft Delete System:**
    * **Soft Delete:** Data is marked as deleted but kept in the DB for audit trails.

## 🔌 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **GET** | `/products` | Get all active products |
| **GET** | `/products/{product_id}` | Get specific product detail |
| **GET** | `/products/search?q=...` | **Search** products by Name |
| **POST** | `/products` | Add a new product |
| **PUT** | `/products/{product_id}` | Update product details (Name, Price) |
| **PATCH** | `/products/{product_id}/stock` | **Adjust Stock** (In/Out logic) |
| **DELETE** | `/products/{product_id}` | **Soft Delete** (Remove from list) |

## 🛠️ How to Run

1.  **Prerequisites:**
    * Ensure PostgreSQL is running.
    * Create a database (e.g., `warehouse_db`).

2.  **Environment Setup:**
    Create a `.env` file in the root directory:
    ```env
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_NAME=warehouse_db
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

## 📦 Data Samples (JSON)

**1. POST /products (Add New Item)**
```json
{
  "product_name": "Mechanical Keyboard",
  "product_price": 750000,
  "product_stock": 50
}
```

**2. PATCH /products/{id}/stock (Stock Opname) Type can be "in" (Add stock) or "out" (Reduce stock).**
```json
{
  "amount": 10,
  "type": "in"
}
```