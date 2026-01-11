# My Golang Backend Learning Projects

A collection of backend projects and experiments built with **Go (Golang)**. This repository documents my learning journey from basic to advance.

## 📂 Projects

| Project | Description | Stack |
| :--- | :--- | :--- |
| **[`/library-api`](./library-api)** | Library management system with borrowing logic. | Go, GORM, PostgreSQL |
| **[`/warehouse-api`](./warehouse-api)** | Warehouse inventory, suppliers, categories, and stock movements management system with audit trails | Go, GORM, PostgreSQL |

## 🚀 How to Run

1.  **Clone the repo:**
    ```bash
    git clone https://github.com/haikalmumtaz233/go-backend-labs.git
    cd go-backend-labs
    ```

2.  **Choose a project and run:**
    ```bash
    cd library-api   # warehouse-api, etc.
    go mod tidy
    go run main.go
    ```

## 🛠️ Requirements

* Go 1.22+
* PostgreSQL (For API projects)
