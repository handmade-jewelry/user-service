-- +goose Up
create table if not exists role
(
    id serial primary key,
    name varchar(20) not null unique
);

-- +goose Down
drop table if exists role;