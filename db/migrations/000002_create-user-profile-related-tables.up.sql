BEGIN;

create type gender as enum ('male', 'female', 'other');

create table if not exists user_profiles
(
    id         bigserial    not null primary key,
    user_id    bigint REFERENCES users (id) UNIQUE,
    first_name varchar(64)  not null,
    last_name  varchar(128) not null,
    gender     gender       not null,
    address    varchar(256),
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now()
);

COMMIT;
