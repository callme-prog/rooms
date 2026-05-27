CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS rooms (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    price       NUMERIC(10,2) NOT NULL,
    capacity    INT          NOT NULL DEFAULT 2,
    image_url   TEXT         NOT NULL DEFAULT ''
);
