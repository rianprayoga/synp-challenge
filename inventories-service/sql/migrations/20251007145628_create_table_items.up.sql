CREATE TABLE items (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    stock INT NOT NULL,
    CONSTRAINT stock_value CHECK (stock >= 0),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null
);