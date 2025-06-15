-- +goose Up
CREATE TABLE task (
    id serial primary key,
    title VARCHAR(250) NOT NULL UNIQUE,
    description text,
    status VARCHAR(20) NOT NULL DEFAULT 'new'
        CHECK (status IN ('new', 'in_progress', 'done', 'todo')),
    author INT NOT NULL,
    operator INT,
    created_at timestamp not null default now(),
    updated_at timestamp,
    due_date timestamp,
    completed_at timestamp
);

COMMENT ON COLUMN task.author IS 'ID автора задачи из auth-сервиса (внешняя БД)';
COMMENT ON COLUMN task.operator IS 'ID оператора задачи из auth-сервиса (внешняя БД)';

-- +goose Down
DROP TABLE person;
