# Easy Commerce Backend

A comprehensive Go-based e-commerce backend API with PostgreSQL, GORM, and RESTful endpoints.

## 🚀 Features

- **User Management**: Registration, authentication, and profile management
- **Product Catalog**: Categories, products, attributes, and variants
- **Shopping Experience**: Wishlist, reviews, and ratings
- **Order Processing**: Cart management and order processing
- **Role-Based Access Control**: Fine-grained permissions system
- **File Management**: Image and file uploads
- **Notifications**: Email and in-app notifications
- **Address Management**: User address book
- **Brand Management**: Product brands and manufacturers
- **Size & Color Management**: Product variations

## 🛠 Prerequisites

- Go 1.24.3 or higher
- PostgreSQL 13+
- Redis (for caching and sessions)
- SMTP server (for email notifications)

## 🚀 Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/easy-commerce-backend.git
   cd easy-commerce-backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   Create a `.env` file in the root directory with the following variables:
   ```env
   # Server
   PORT=8080
   ENV=development
   
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_secure_password
   DB_NAME=easy_commerce
   DB_SSL_MODE=disable
   
   # JWT
   JWT_SECRET=your_jwt_secret
   JWT_EXPIRATION=24h
   
   # Redis
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   
   # SMTP (for emails)
   SMTP_HOST=smtp.example.com
   SMTP_PORT=587
   SMTP_USERNAME=your_email@example.com
   SMTP_PASSWORD=your_email_password
   SMTP_FROM=noreply@easycommerce.com
   ```

4. **Run migrations**
   ```bash
   go run cmd/migrate/main.go
   ```

5. **Start the server**
   ```bash
   go run cmd/main.go
   ```

   The API will be available at `http://localhost:8080`

## 📚 API Documentation

For detailed API documentation, please refer to the [API Documentation](api-doc/README.md) or import the Postman collection from the `postman` directory.

## 🏗 Project Structure

```
backend/
├── cmd/
│   ├── main.go              # Application entry point
│   └── migrate/             # Database migration scripts
├── db/
│   └── migrations/          # SQL migration files
├── internal/
│   ├── address/             # User address management
│   ├── attribute_option/    # Product attribute options
│   ├── attribute_type/      # Product attribute types
│   ├── auth/                # Authentication and authorization
│   ├── brand/               # Product brands
│   ├── category/            # Product categories
│   ├── color/               # Product colors
│   ├── file_storage/        # File upload and storage
│   ├── mail/                # Email services
│   ├── notification/        # Notification system
│   ├── permission/          # Permission management
│   ├── product_stats/       # Product statistics
│   ├── review/              # Product reviews and ratings
│   ├── role/                # User roles
│   ├── size_category/       # Size categories
│   ├── size_option/         # Size options
│   ├── system/              # System configurations
│   ├── user/                # User management
│   └── wishlist/            # User wishlists
├── pkg/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection and models
│   ├── logger/              # Logging utilities
│   ├── middleware/          # HTTP middlewares
│   ├── response/            # API response helpers
│   └── utils/               # Utility functions
└── route/                   # API route definitions
```

## 🔒 Authentication

Most endpoints require authentication. Include the JWT token in the `Authorization` header:

```
Authorization: Bearer your_jwt_token_here
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Server Shutdown Gracefully
lsof -ti:8080 | xargs kill && sleep 3 || lsof -ti:8080 | xargs kill -9
