-- +goose Up
CREATE TABLE task_watchers (
    task_id INT NOT NULL REFERENCES task(id) ON DELETE CASCADE,
    watcher INT not null,
    PRIMARY KEY (task_id, watcher)
);

-- +goose Down
DROP TABLE task_watchers;
