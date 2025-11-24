# 🛒 Easy Commerce Backend

[![Go Version](https://img.shields.io/badge/Go-1.25.0-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GORM](https://img.shields.io/badge/ORM-GORM-red?style=flat)](https://gorm.io/)

A comprehensive, production-ready Go-based e-commerce backend API built with modern architecture patterns, featuring PostgreSQL database, GORM ORM, RESTful endpoints, and cloud storage integration.

---

## ✨ Features

### 🔐 Authentication & Authorization
- **JWT-based Authentication** with secure token management
- **OAuth Integration** (Google, Facebook)
- **Role-Based Access Control (RBAC)** with fine-grained permissions
- **Super Admin** initialization on first run
- Social login with default password fallback

### 👥 User Management
- User registration and profile management
- Address book management
- User preferences and settings
- Email verification and password reset

### 🛍️ Product Management
- **Categories** with hierarchical structure
- **Brands** and manufacturer management
- **Product Attributes** (types and options)
- **Size & Color Variants** with categories
- **Product Statistics** tracking
- **Reviews & Ratings** system

### 🛒 Shopping Experience
- **Wishlist** functionality
- Product search and filtering
- Advanced sorting options
- Product recommendations

### 📦 Order & Inventory
- Cart management
- Order processing workflow
- Inventory tracking
- Product statistics

### 📁 File Management
- **Cloudflare R2** integration for object storage
- **AWS S3-compatible** storage provider
- Image upload and optimization
- Secure file access control

### 📧 Communication
- **Email Service** with SMTP integration
- **Push Notifications** via FCM (Firebase Cloud Messaging)
- In-app notification system
- Template-based email system

### 🎨 Content Management
- **CMS Module** for dynamic content
- Contact information management
- System configuration management

### 🔧 Technical Features
- **UUID Support** (sequence or random generation)
- **Graceful Shutdown** handling
- **Structured Logging** with zerolog
- **Database Migrations** with version control
- **CORS** middleware
- **Recovery** middleware for panic handling
- **Request Logging** middleware
- **Environment-based Configuration** (dev, stage, prod)

---

## 🛠 Tech Stack

| Category | Technology |
|----------|-----------|
| **Language** | Go 1.25.0 |
| **Database** | PostgreSQL 13+ |
| **ORM** | GORM v1.31.1 |
| **Authentication** | JWT (golang-jwt/jwt v5.3.0) |
| **Cloud Storage** | AWS S3 SDK v2 / Cloudflare R2 |
| **Logging** | zerolog v1.34.0 |
| **Email** | go-mail v0.7.2 |
| **Environment** | godotenv v1.5.1 |
| **OAuth** | golang.org/x/oauth2 v0.31.0 |
| **UUID** | google/uuid v1.6.0 |

---

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.25.0 or higher ([Download](https://go.dev/dl/))
- **PostgreSQL** 13+ ([Download](https://www.postgresql.org/download/))
- **Git** ([Download](https://git-scm.com/downloads))
- **Cloudflare R2** or **AWS S3** account (for file storage)
- **SMTP Server** (optional, for email notifications)
- **FCM Server Key** (optional, for push notifications)

---

## 🚀 Quick Start

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/zakaria5729/golang-ecommerce.git
cd golang-ecommerce
```

### 2️⃣ Install Dependencies

```bash
go mod tidy
```

### 3️⃣ Environment Configuration

The project supports multiple environments. Choose one based on your needs:

- **Development**: `.env.dev`
- **Local**: `.env.local`
- **Staging**: `.env.stage`
- **Production**: `.env.prod`

Copy the appropriate environment file or create your own `.env`:

```bash
cp .env.dev .env
```

#### Environment Variables Reference

```env
# Server Configuration
PORT=8080
HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=eascos
DB_SSL_MODE=disable
DB_SHOW_CONSOLE_LOG=true

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_change_this_in_production

# Super Admin (Created on first run)
SUPER_ADMIN_EMAIL=admin@admin.com
SUPER_ADMIN_PASSWORD=admin@admin

# Application URLs
DOMAIN_URL=http://localhost:${PORT}

# Firebase Cloud Messaging (Optional)
FCM_SERVER_KEY=your_fcm_server_key
FCM_URL=https://fcm.googleapis.com/fcm/send

# OAuth Configuration (Optional)
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
FACEBOOK_APP_ID=your_facebook_app_id

# Object Storage (Cloudflare R2 / AWS S3)
OBJ_STORE_REGION=auto
OBJ_STORE_BUCKET_NAME=your_bucket_name
OBJ_STORE_ACCOUNT_ID=your_account_id
OBJ_STORE_ACCESS_KEY_ID=your_access_key_id
OBJ_STORE_ACCESS_KEY_SECRET=your_access_key_secret
OBJ_STORE_PUBLIC_DOMAIN=your_public_domain.r2.dev

# Application Settings
APP_HEALTH_CHECK_TOKEN=your_health_check_token
SOCIAL_LOGIN_DEFAULT_PASSWORD=default_social_password

# UUID Generation Type: sequence / random
UUID_TYPE=sequence

# Logging Configuration: always / error
WRITE_LOG_WHEN=always

# Database Stats Logging: sql / sql_analyze
DB_STATS_LOG_TYPE=sql_analyze
```

### 4️⃣ Database Setup

Create a PostgreSQL database:

```bash
createdb eascos
```

Or using SQL:

```sql
CREATE DATABASE eascos;
```

### 5️⃣ Run Database Migrations

The application automatically runs migrations on startup. Alternatively, you can run migrations manually:

```bash
go run cmd/main.go
```

### 6️⃣ Start the Server

```bash
go run cmd/main.go
```

You should see:

```
Server started successfully on port 8080
```

The API will be available at: **http://localhost:8080**

---

## 🏗️ Project Structure

```
golang-ecommerce/
├── cmd/
│   ├── main.go                    # Application entry point
│   └── uuid_migrate.go            # UUID migration utility
│
├── db/
│   ├── db.go                      # Database initialization
│   ├── migration2/                # SQL migration files
│   │   ├── create_all_tables.sql
│   │   ├── uuid_migration.sql
│   │   └── action_date_default_null.sql
│   └── migrations/                # Migration management
│       ├── migrate.go
│       └── migrate2.go
│
├── internal/                      # Internal application modules
│   ├── address/                   # User address management
│   ├── attribute_option/          # Product attribute options
│   ├── attribute_type/            # Product attribute types
│   ├── auth/                      # Authentication & authorization
│   ├── brand/                     # Brand management
│   ├── category/                  # Category management
│   ├── cms/                       # Content management system
│   ├── color/                     # Color variants
│   ├── contact_info/              # Contact information
│   ├── file_storage/              # File upload & storage
│   │   └── provider/              # Storage providers (R2, S3)
│   ├── mail/                      # Email service
│   ├── notification/              # Notification system
│   ├── permission/                # Permission management
│   ├── product_stats/             # Product statistics
│   ├── review/                    # Product reviews & ratings
│   ├── role/                      # Role management
│   ├── size_category/             # Size categories
│   ├── size_option/               # Size options
│   ├── system/                    # System configuration
│   │   └── model/                 # System models
│   ├── user/                      # User management
│   └── wishlist/                  # Wishlist functionality
│
├── pkg/                           # Shared packages
│   ├── base/                      # Base entities & repositories
│   ├── config/                    # Configuration management
│   ├── constants/                 # Application constants
│   ├── contextutil/               # Context utilities
│   ├── data_loader/               # Initial data loader
│   ├── logger/                    # Logging utilities
│   ├── middleware/                # HTTP middlewares
│   ├── option/                    # Query & sort options
│   ├── router/                    # Router utilities
│   ├── tokenutil/                 # JWT token utilities
│   └── utils/                     # General utilities
│
├── route/                         # API route definitions
│   ├── routes.go                  # Route registration
│   ├── auth_route.go
│   ├── user_route.go
│   ├── role_route.go
│   ├── permission_route.go
│   ├── category_route.go
│   ├── brand_route.go
│   ├── product_stats_route.go
│   ├── review_route.go
│   ├── wishlist_route.go
│   ├── file_storage_route.go
│   ├── system_route.go
│   └── ... (other routes)
│
├── logs/                          # Application logs
├── templates/                     # Email templates
├── .env.dev                       # Development environment
├── .env.local                     # Local environment
├── .env.stage                     # Staging environment
├── .env.prod                      # Production environment
├── .gitignore
├── go.mod                         # Go module dependencies
├── go.sum                         # Dependency checksums
└── README.md                      # This file
```

---

## 📡 API Endpoints Overview

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh JWT token
- `POST /api/auth/google` - Google OAuth login
- `POST /api/auth/facebook` - Facebook OAuth login

### User Management
- `GET /api/users` - List users (Admin)
- `GET /api/users/:id` - Get user details
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Roles & Permissions
- `GET /api/roles` - List roles
- `POST /api/roles` - Create role
- `PUT /api/roles/:id` - Update role
- `GET /api/permissions` - List permissions

### Products & Categories
- `GET /api/categories` - List categories
- `POST /api/categories` - Create category
- `GET /api/brands` - List brands
- `GET /api/colors` - List colors
- `GET /api/sizes` - List sizes

### Reviews & Wishlist
- `GET /api/reviews` - List reviews
- `POST /api/reviews` - Create review
- `GET /api/wishlist` - Get user wishlist
- `POST /api/wishlist` - Add to wishlist

### File Storage
- `POST /api/files/upload` - Upload file
- `GET /api/files/:id` - Get file details

### System
- `GET /api/system/health` - Health check
- `GET /api/system/config` - System configuration

> **Note**: Most endpoints require JWT authentication. Include the token in the `Authorization` header:
> ```
> Authorization: Bearer <your_jwt_token>
> ```

---

## 🔒 Authentication Flow

1. **Register** a new user via `/api/auth/register`
2. **Login** with credentials to receive a JWT token
3. **Include** the JWT token in subsequent requests
4. **Refresh** the token before expiration

### Example Login Request

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your_password"
  }'
```

### Example Authenticated Request

```bash
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 🧪 Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o easy-commerce cmd/main.go
```

### Running the Binary

```bash
./easy-commerce
```

---

## 🚢 Deployment

### Docker Deployment (Recommended)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.25.0-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env.prod .env
EXPOSE 8080
CMD ["./main"]
```

Build and run:

```bash
docker build -t easy-commerce .
docker run -p 8080:8080 easy-commerce
```

### Manual Deployment

1. Build the binary for your target platform
2. Copy the binary and `.env.prod` to your server
3. Set up PostgreSQL database
4. Run migrations
5. Start the application with a process manager (systemd, supervisor, pm2)

---

## 🛠️ Utilities

### Graceful Server Shutdown

The application handles graceful shutdown automatically. To manually stop the server:

```bash
# Graceful shutdown
lsof -ti:8080 | xargs kill && sleep 3 || lsof -ti:8080 | xargs kill -9
```

### UUID Migration

If you need to migrate from integer IDs to UUIDs:

```bash
go run cmd/uuid_migrate.go
```

---

## 📊 Logging

Logs are stored in the `logs/` directory with daily rotation:

- **Format**: JSON structured logging
- **Levels**: Debug, Info, Warn, Error
- **Rotation**: Daily
- **Configuration**: Set `WRITE_LOG_WHEN` to `always` or `error`

---

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/AmazingFeature`)
3. **Commit** your changes (`git commit -m 'Add some AmazingFeature'`)
4. **Push** to the branch (`git push origin feature/AmazingFeature`)
5. **Open** a Pull Request

### Code Style

- Follow Go best practices and conventions
- Use `gofmt` for code formatting
- Write meaningful commit messages
- Add tests for new features

---

## 📝 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 👨‍💻 Author

**Zakaria**
- GitHub: [@zakaria5729](https://github.com/zakaria5729)
- Repository: [golang-ecommerce](https://github.com/zakaria5729/golang-ecommerce)

---

## 🙏 Acknowledgments

- [GORM](https://gorm.io/) - The fantastic ORM library
- [zerolog](https://github.com/rs/zerolog) - Zero allocation JSON logger
- [golang-jwt](https://github.com/golang-jwt/jwt) - JWT implementation
- [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/) - Cloud storage integration

---

## 📞 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/zakaria5729/golang-ecommerce/issues) page
2. Create a new issue with detailed information
3. Contact the maintainer

---

**Happy Coding! 🚀**
