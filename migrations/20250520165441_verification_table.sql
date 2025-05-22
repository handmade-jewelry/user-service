-- +goose Up
create table if not exists verification
(
    token text primary key,
    user_id int not null,
    is_used boolean default false,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    expired_at timestamp not null
);

-- +goose Down
drop table if exists verification;
