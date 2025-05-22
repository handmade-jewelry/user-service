-- +goose Up
create table if not exists user_role
(
    user_id int not null,
    role_id int not null,
    assigned_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp,
    constraint user_role_unique unique (user_id, role_id)
);

-- +goose Down
drop table if exists user_role;
