-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
                                user_id UUID PRIMARY KEY,
                                token_hash TEXT NOT NULL,
                                client_ip TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
