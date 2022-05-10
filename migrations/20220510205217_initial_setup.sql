-- +goose Up
-- +goose StatementBegin
BEGIN;


COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
