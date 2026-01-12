# My Golang Backend Learning Projects

A collection of backend projects and experiments built with **Go (Golang)**. This repository documents my learning journey from basic to advance.

## 📂 Projects

| Project | Description | Stack |
| :--- | :--- | :--- |
| **[`/library-api`](./library-api)** | Library management system with borrowing logic. | Go, GORM, PostgreSQL |
| **[`/warehouse-api`](./warehouse-api)** | Warehouse inventory, suppliers, categories, and stock movements management system with audit trails | Go, GORM, PostgreSQL |
| **[`/warehouse-gin`](./warehouse-gin)** | Warehouse inventory, suppliers, categories, and stock movements management system with audit trails | Go, Gin, GORM, PostgreSQL |
| **[`/eventix`](./eventix)** | A scalable, production-ready backend for an event ticketing platform | Go, Gin, GORM, PostgreSQL |


## 🚀 How to Run

### Option 1: Manual Run (All Projects)

1.  **Clone the repo:**
    ```bash
    git clone https://github.com/haikalmumtaz233/go-backend-labs.git
    cd go-backend-labs
    ```

2.  **Choose a project and run:**
    ```bash
    cd library-api   # or warehouse-api, eventix
    go mod tidy
    go run main.go
    ```

### Option 2: Docker Compose (Eventix Only)

```bash
cd eventix
docker-compose up --build
```

## 🛠️ Requirements

* Go 1.25+
* PostgreSQL (For API projects)
* Docker & Docker Compose (Optional, for Eventix)

---

## License

This project is licensed under the MIT License.
