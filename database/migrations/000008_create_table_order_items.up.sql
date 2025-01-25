CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);

CREATE INDEX idx_orders_items_customer_id ON order_items(order_id);
CREATE INDEX idx_orders_items_product_id ON order_items(product_id);
CREATE INDEX idx_orders_items_price ON order_items(price);
