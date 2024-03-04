-- +goose Up
-- +goose StatementBegin
create table if not exists messages
(
    id                  serial primary key,
    text                text not null,
    role                text not null,
    message_id          bigint not null,
    chat_id             bigint not null,
    context_uuid        uuid not null
);

create index idx_messages on messages (context_uuid, message_id);

create table if not exists users_whitelist
(
    user_id             bigint unique not null,
    state               bool not null,
    description         text not null
);

create index idx_users_whitelist on users_whitelist (user_id);

create table if not exists chats_whitelist
(
    chat_id             bigint unique not null,
    state               bool not null,
    description         text not null
);

create index idx_chats_whitelist on chats_whitelist (chat_id);

create table if not exists ai_presets
(
    id                  serial primary key,
    chat_id             bigint not null,
    text                text not null,
    tag                 text not null
);

create index idx_ai_presets on ai_presets (chat_id, tag);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
