BEGIN;

create type user_status as ENUM ('active', 'inactive', 'deleted');
create type user_roles as ENUM ('admin', 'moderator', 'merchant', 'user');

create table if not exists users
(
    id          bigserial    not null primary key,
    user_id     uuid         not null default uuid_generate_v4() UNIQUE,
    username    citext       not null unique,
    email       citext       not null unique,
    status      user_status  not null default 'active',
    role        user_roles   not null default 'user',
    hashed_pass varchar(128) not null,
    created_at  timestamptz  not null default now(),
    updated_at  timestamptz  not null default now()
);

ALTER TABLE users
    ADD CONSTRAINT username_length CHECK (length(username) <= 64);
ALTER TABLE users
    ADD CONSTRAINT email_length CHECK (length(email) <= 128);

COMMIT;
