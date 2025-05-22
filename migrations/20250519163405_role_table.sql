-- +goose Up
create table if not exists role
(
    id serial primary key,
    name varchar(20) not null unique,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp
);

-- +goose Down
drop table if exists role;