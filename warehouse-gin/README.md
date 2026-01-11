# 🍸 Warehouse Management API (Gin Version)

This project is a refactored version of my previous Warehouse API.
It demonstrates the migration from Go's standard `net/http` library to **Gin Gonic**.

## 📐 Project Architecture

The application is structured into 4 distinct layers:

```text
wms-clean-arch/
├── entity/       # 🟢 Domain Models
├── repository/   # 🔵 Data Access Layer
├── service/      # 🟡 Business Logic Layer
├── handler/      # 🔴 HTTP Transport Layer
└── main.go       # ⚙️ Dependency Injection & Wiring.
```

## ⚡ Key Features

* **Clean Architecture:**
  * **Dependency Injection:** Layers are connected in `main.go` (Handler -> Service -> Repository).
  * **Interface-Oriented:** Services and Repositories rely on interfaces, making the code testable and flexible.
* **Business Logic Isolation:** Logic like "Stock Checks" or "Validation" lives strictly in the Service layer, not in Handlers.
* **Database:** PostgreSQL with GORM.
* **Atomic Transactions:** Stock adjustments (In/Out) + Audit Logging are wrapped in DB Transactions (`tx`) within the Repository layer.
* **Audit Trail:** Automatically records stock mutation history.

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