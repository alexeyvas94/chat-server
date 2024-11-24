-- +goose Up
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE chat_users (
    chat_id INT REFERENCES chats(id) ON DELETE CASCADE,
    user_name VARCHAR(30) NOT NULL,
    PRIMARY KEY (chat_id, user_name)
);

CREATE TABLE message (
    id SERIAL PRIMARY KEY,
    from_user VARCHAR(30) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS chats, message, chat_users;