CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    ip_uid TEXT NOT NULL,
    id_on_stripe TEXT NOT NULL
);

CREATE TABLE stripe_products (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    id_on_stripe TEXT NOT NULL
);

CREATE TABLE stripe_prices (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    id_on_stripe TEXT NOT NULL,
    type TEXT NOT NULL,
    stripe_product_id BIGINT NOT NULL
);

-- typeがrecurringの場合のみ紐づく
CREATE TABLE stripe_recurrings (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    interval TEXT NOT NULL,
    interval_count SMALLINT NOT NULL,
    stripe_price_id BIGINT NOT NULL
);

CREATE TABLE stripe_order_items (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    stripe_price_id BIGINT NOT NULL
);

-- add constraint

ALTER TABLE stripe_prices
ADD CONSTRAINT fk_stripe_prices_stripe_products
FOREIGN KEY (stripe_product_id)
REFERENCES stripe_products (id);
