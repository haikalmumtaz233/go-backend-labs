# 📦 Warehouse Management API

A robust RESTful API to manage warehouse inventory, suppliers, categories, and stock movements with audit trails.

## 🚀 Key Features

* **Relational Database:**
    * **Products:** Linked with **Categories** and **Suppliers**.
    * **Associations:** Uses GORM Preloading for fetching related data.
* **Audit Trail System:**
    * **Stock History:** Records every stock change (In/Out) in a `stock_mutations` table.
    * **Atomic Transactions:** Uses Database Transactions (`tx`) to ensure stock updates and log recording happen simultaneously (ACID compliant).
* **Security:**
    * **Anti-SQL Injection:** Uses GORM parameterized queries.
    * **Environment Variables:** Sensitive data managed via `.env`.
* **Smart Stock Management:**
    * **Validation:** Prevents negative stock during outgoing adjustments.

## 🔌 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **GET** | `/products` | Get all products (with Category & Supplier) |
| **GET** | `/products/{id}` | Get specific product detail |
| **GET** | `/products/search?q=...` | **Search** products by Name |
| **GET** | `/products/{id}/history` | **Audit Log:** View stock mutation history |
| **POST** | `/products` | Add a new product |
| **PUT** | `/products/{id}` | Update product details |
| **PATCH** | `/products/{id}/stock` | **Adjust Stock** (In/Out with Transaction) |
| **DELETE** | `/products/{id}` | **Soft Delete** product |
| **GET** | `/categories` | Get all categories |
| **POST** | `/categories` | Add a new category |
| **GET** | `/suppliers` | Get all suppliers |
| **POST** | `/suppliers` | Add a new supplier |

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
    go run .
    ```
    Server runs at `http://localhost:8080`.

## 📦 Data Samples (JSON)

**1. POST /products (Add New Item)**
*Requires existing Category ID and Supplier ID*
```json
{
  "product_name": "Gaming Laptop",
  "product_price": 15000000,
  "product_stock": 10,
  "category_id": 1,
  "supplier_id": 1
}
```
**2. PATCH /products/{id}/stock (Stock Opname) Recorded in Audit Log**
```json
{
  "amount": 5,
  "type": "in"
}
```