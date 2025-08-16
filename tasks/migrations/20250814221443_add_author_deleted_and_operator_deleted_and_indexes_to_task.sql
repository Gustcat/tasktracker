-- +goose Up
ALTER TABLE task
    ADD COLUMN author_deleted BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN operator_deleted BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX idx_task_author_deleted ON task(author_deleted);
CREATE INDEX idx_task_operator_deleted ON task(operator_deleted);

-- +goose Down
DROP INDEX IF EXISTS idx_task_author_deleted;
DROP INDEX IF EXISTS idx_task_operator_deleted;

ALTER TABLE task
    DROP COLUMN IF EXISTS author_deleted,
    DROP COLUMN IF EXISTS operator_deleted;