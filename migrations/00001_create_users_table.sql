-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    user_id uuid PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    phone_number varchar(15) UNIQUE,
    PASSWORD TEXT NOT NULL,
    user_type varchar(10) CHECK (user_type IN ('admin', 'user')) NOT NULL,
    is_active boolean NOT NULL DEFAULT TRUE,
    created_at timestamp WITH time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp WITH time zone NOT NULL DEFAULT NOW(),
    deleted_at timestamp WITH time zone
);

CREATE INDEX idx_users_email ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
