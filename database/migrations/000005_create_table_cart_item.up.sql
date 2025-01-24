CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INTEGER NOT NULL REFERENCES Cart(id),
    product_id INTEGER NOT NULL REFERENCES Product(id),
    quantity INTEGER NOT NULL
);

CREATE INDEX idx_cart_item_cart_id ON cart_items(cart_id);
CREATE INDEX idx_cart_item_product_id ON cart_items(product_id);