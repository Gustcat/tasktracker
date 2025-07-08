-- +goose Up
CREATE TABLE task_watchers (
    id serial primary key,
    task_id INT NOT NULL REFERENCES task(id) ON DELETE CASCADE,
    watcher INT NOT NULL,
    UNIQUE(task_id, watcher)
);

COMMENT ON COLUMN task_watchers.watcher IS 'ID наблюдателя задачи из auth-сервиса (внешняя БД)';

-- +goose Down
DROP TABLE task_watchers;
