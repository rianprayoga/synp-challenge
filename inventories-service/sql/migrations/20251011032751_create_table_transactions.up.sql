CREATE TABLE transactions (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE,
    item_id UUID NOT NULL REFERENCES items (id),
    qty INT NOT NULL,
    released BOOLEAN DEFAULT false,
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null
);