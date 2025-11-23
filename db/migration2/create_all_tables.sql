ALTER DATABASE easy_commerce SET timezone TO 'UTC';
SHOW timezone;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    verified BOOLEAN DEFAULT FALSE,
    banned BOOLEAN DEFAULT FALSE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP,
    verification_token VARCHAR(255),
    verification_expires TIMESTAMP,
    refresh_token_expires TIMESTAMP,
    refresh_token VARCHAR(255),
    path_key VARCHAR(255),
    purchase_count INTEGER DEFAULT 0,
    total_spent INTEGER DEFAULT 0,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    role_type VARCHAR(50) NOT NULL,
    role_name VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(255),
    group_name VARCHAR(100) DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    is_active BOOLEAN DEFAULT true,
    parent_id INTEGER,
    priority INTEGER DEFAULT 0,
    title VARCHAR(120) NOT NULL,
    sub_title VARCHAR(255),
    path_key VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    is_default BOOLEAN DEFAULT FALSE,
    user_id INTEGER NOT NULL,
    address_type VARCHAR(20),
    zip_code VARCHAR(20),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    country VARCHAR(100) NOT NULL,
    street TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description VARCHAR(2000),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS colors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(2000),
    price DECIMAL(10,2) NOT NULL,
    stock_quantity INTEGER DEFAULT 0,
    sku VARCHAR(100) UNIQUE,
    is_active BOOLEAN DEFAULT true,
    category_id INTEGER NOT NULL,
    brand_id INTEGER,
    color_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE SET NULL,
    FOREIGN KEY (color_id) REFERENCES colors(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS product_stats (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    view_count INTEGER NOT NULL DEFAULT 0,
    add_to_cart_count INTEGER NOT NULL DEFAULT 0,
    remove_from_cart_count INTEGER NOT NULL DEFAULT 0,
    wishlist_count INTEGER NOT NULL DEFAULT 0,
    remove_from_wishlist_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS browsing_histories (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    viewed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    product_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS wishlists (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE(user_id, product_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS attribute_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS attribute_options (
    id SERIAL PRIMARY KEY,
    attribute_type_id INTEGER NOT NULL,
    attribute_option_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (attribute_type_id) REFERENCES attribute_types(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS size_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS size_options (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    sort_order INTEGER DEFAULT NULL,
    size_category_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (size_category_id) REFERENCES size_categories(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS contact_infos (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    name VARCHAR(255) DEFAULT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL,
    phone VARCHAR(255) DEFAULT NULL,
    message TEXT DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS cms_pages (
    id SERIAL PRIMARY KEY,
    tag VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS cms_sections (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) DEFAULT NULL,
    link VARCHAR(255) DEFAULT NULL,
    content TEXT NOT NULL,
    page_id INT NOT NULL,
    image_urls JSONB DEFAULT '[]'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by INT DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (page_id) REFERENCES cms_pages(id) ON DELETE RESTRICT
);