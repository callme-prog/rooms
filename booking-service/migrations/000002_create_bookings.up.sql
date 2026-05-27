CREATE TABLE IF NOT EXISTS bookings (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL,
    room_id    UUID        NOT NULL REFERENCES rooms(id),
    check_in   TIMESTAMPTZ NOT NULL,
    check_out  TIMESTAMPTZ NOT NULL,
    status     VARCHAR(30) NOT NULL DEFAULT 'pending',
    comment    TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_dates CHECK (check_out > check_in)
);
