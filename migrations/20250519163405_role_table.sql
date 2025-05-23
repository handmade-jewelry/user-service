-- +goose Up
create table if not exists role
(
    id serial primary key,
    name varchar(20) not null unique,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp
);

INSERT INTO role (name)
VALUES
    ('SELLER'),
    ('CUSTOMER'),
    ('ADMIN')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
drop table if exists role;