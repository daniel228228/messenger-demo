create extension if not exists "uuid-ossp";

create table "user"
(
    id         uuid default gen_random_uuid() not null
        constraint user_pk
            primary key,
    username   text,
    first_name text,
    last_name  text
);