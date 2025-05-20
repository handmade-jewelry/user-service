-- +goose Up
create table if not exists users
(
    id serial primary key,
    email varchar(50) not null unique,
    email_verified boolean default false,
    password_hash text not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp
);

-- +goose Down
drop table if exists users;
