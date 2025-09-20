CREATE TABLE IF NOT EXISTS product_stats (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    view_count INT NOT NULL DEFAULT 0,
    add_to_cart_count INT NOT NULL DEFAULT 0,
    remove_from_cart_count INT NOT NULL DEFAULT 0,
    wishlist_count INT NOT NULL DEFAULT 0,
    remove_from_wishlist_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- add foreign key to product table
-- ALTER TABLE product_stats ADD CONSTRAINT fk_product_stats_product_id FOREIGN KEY (product_id) REFERENCES products(id);

-- CREATE INDEX idx_product_stats_product_id ON product_stats(product_id);

INSERT INTO product_stats (product_id, view_count, add_to_cart_count, remove_from_cart_count, wishlist_count, remove_from_wishlist_count) VALUES
(1, 0, 0, 0, 0, 0),
(2, 0, 0, 0, 0, 0),
(3, 0, 0, 0, 0, 0);