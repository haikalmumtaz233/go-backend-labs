# 🍸 Warehouse Management API (Gin Version)

This project is a refactored version of my previous Warehouse API.
It demonstrates the migration from Go's standard `net/http` library to **Gin Gonic**.

## 🚀 Upgrade Highlights

Comparing this version to the previous standard library version:

* **Framework:** Switched to **Gin Gonic** for cleaner syntax and speed.
* **Routing:** Implemented **Group Routing** (`/api`) and cleaner parameter handling (`:id`).
* **Request Parsing:** Replaced manual `json.Decode` with Gin's `ShouldBindJSON` for automatic binding and validation.
* **Response:** Using `c.JSON` for standardized JSON responses.

## ⚡ Key Features

* **Database:** PostgreSQL with **GORM**.
* **Complex Relations:** One-to-Many (Category -> Products, Supplier -> Products).
* **Atomic Transactions:** Stock adjustments (In/Out) are wrapped in DB Transactions to ensure data integrity.
* **Audit Logging:** Automatically records every stock movement into a history table.
* **Search:** Pattern matching search using SQL `ILIKE`.

## 🔌 API Endpoints

**Base URL:** `http://localhost:8080/api`

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **GET** | `/products` | List all products (with Category & Supplier) |
| **GET** | `/products/search?q=...` | **Search** products by Name |
| **GET** | `/products/:id` | Get product detail |
| **GET** | `/products/:id/history` | View **Stock Mutation Logs** |
| **POST** | `/products` | Create new product |
| **PUT** | `/products/:id` | Update product info |
| **PATCH** | `/products/:id/stock` | **Adjust Stock** (In/Out logic) |
| **DELETE** | `/products/:id` | Soft delete product |
| **GET** | `/categories` | List categories |
| **POST** | `/categories` | Create category |
| **GET** | `/suppliers` | List suppliers |
| **POST** | `/suppliers` | Create supplier |

## 🛠️ How to Run

1.  **Clone & Install Dependencies**
    ```bash
    go mod tidy
    ```

2.  **Environment Setup**
    Create a `.env` file:
    ```env
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_NAME=warehouse_db
    DB_PORT=5432
    ```

3.  **Run Server**
    ```bash
    go run .
    ```

## 📦 Data Samples (JSON)

**1. Add New Item**
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
**2. Adjust Stock (PATCH)**
```json
{
  "amount": 5,
  "type": "in"
}
```