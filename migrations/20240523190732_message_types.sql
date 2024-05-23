-- +goose Up
-- +goose StatementBegin

alter table messages add column type text default 'text';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
