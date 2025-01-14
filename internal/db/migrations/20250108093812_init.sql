-- +goose Up
CREATE TABLE todos(
    id                  BIGSERIAL       PRIMARY KEY,
    title               VARCHAR(20)     NOT NULL,
    description         TEXT            NOT NULL,
    completed           BOOLEAN         NOT NULL        DEFAULT FALSE,
    completed_at        TIMESTAMPTZ     NULL,
    created_at          TIMESTAMPTZ     NOT NULL        DEFAULT clock_timestamp(),
    updated_at          TIMESTAMPTZ     NOT NULL        DEFAULT clock_timestamp()
);

CREATE TABLE todos_comments(
    id                  BIGSERIAL       PRIMARY KEY,
    comment             TEXT            NOT NULL,
    todo_id             INT             NOT NULL,
    created_at          TIMESTAMPTZ     NOT NULL        DEFAULT clock_timestamp()
);

ALTER TABLE todos_comments ADD CONSTRAINT todo_comments_todo_id FOREIGN KEY (todo_id) REFERENCES todos (id);

-- +goose Down
ALTER TABLE todos_comments DROP CONSTRAINT todo_comments_todo_id;
DROP TABLE IF EXISTS todos_comments;
DROP TABLE IF EXISTS todos;