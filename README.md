# Easy Commerce Backend

A Go-based e-commerce backend API with PostgreSQL and GORM.

## Features

- PostgreSQL database with GORM ORM
- Category management API
- RESTful endpoints with query parameter filtering
- Auto-migration and database setup

## Prerequisites

- Go 1.24.3 or higher
- PostgreSQL database

## Setup

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Set up PostgreSQL database:**
   - Create a database named `easy_commerce`
   - Update environment variables or use defaults

3. **Environment Variables:**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=your_password
   export DB_NAME=easy_commerce
   export DB_SSL_MODE=disable
   export PORT=8080
   ```

4. **Run migrations (required for first setup):**
   ```bash
   go run cmd/migrate/main.go
   ```

5. **Run the application:**
   ```bash
   go run cmd/main.go
   ```

## API Endpoints

### Categories

#### Get All Categories
```
GET /api/v1/categories
```

**Query Parameters:**
- `include` - Comma-separated list of fields to include (sub_title, image_url, is_active, parent_id)
- `parent_id` - Filter by parent category ID (use "null" for root categories)
- `is_active` - Filter by active status (true/false)

**Examples:**
```bash
# Get all categories
curl http://localhost:8080/api/v1/categories

# Get categories with additional fields
curl "http://localhost:8080/api/v1/categories?include=sub_title,image_url,is_active"

# Get active categories only
curl "http://localhost:8080/api/v1/categories?is_active=true"

# Get root categories (no parent)
curl "http://localhost:8080/api/v1/categories?parent_id=null"
```

#### Get Single Category
```
GET /api/v1/categories/{id}
```

**Query Parameters:**
- `include` - Comma-separated list of fields to include

**Example:**
```bash
curl "http://localhost:8080/api/v1/categories/1?include=sub_title,image_url"
```

## Response Format

### Basic Category Response
```json
{
    "id": 1,
    "title": "Electronics",
    "created_at": "2024-01-15 10:30:45.123456",
    "updated_at": "2024-01-15 10:30:45.123456"
}
```

### Extended Category Response (with include parameter)
```json
{
    "id": 1,
    "title": "Electronics",
    "sub_title": "Electronic devices and gadgets",
    "image_url": "https://example.com/electronics.jpg",
    "parent_id": null,
    "is_active": true,
    "created_at": "2024-01-15 10:30:45.123456",
    "updated_at": "2024-01-15 10:30:45.123456"
}
```

## Database Schema

The application uses SQL migrations to create the database schema. The main table is:

### Categories Table
- `id` - Primary key (auto-increment)
- `title` - Category title (required)
- `sub_title` - Optional subtitle
- `image_url` - Optional image URL
- `parent_id` - Optional parent category ID (self-referencing)
- `is_active` - Active status (default: true)
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp
- `updated_at` - Last update timestamp

## Project Structure

```
backend/
├── cmd/
│   └── main.go              # Application entry point
├── db/
│   ├── db.go                # Database connection
│   └── migrations/          # SQL migration files
├── internal/
│   ├── feature/
│   │   └── category/
│   │       ├── category_model.go    # Category model
│   │       ├── category_repository.go # Database operations
│   │       └── category_usecase.go   # Business logic
│   ├── handler/
│   │   └── category_handler.go # HTTP handlers
│   └── route/
│       └── routes.go        # Route definitions
├── pkg/
│   ├── config/
│   │   └── env_config.go    # Configuration management
│   └── response/
│       └── response.go      # HTTP response helpers
└── go.mod                   # Go module file
```
