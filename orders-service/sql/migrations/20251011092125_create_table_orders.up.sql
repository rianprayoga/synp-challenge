CREATE TABLE orders (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL,
    qty INT NOT NULL,
    isConfirmed BOOLEAN NOT NULL,
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null
);