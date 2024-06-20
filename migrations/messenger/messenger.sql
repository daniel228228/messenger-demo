create extension if not exists "uuid-ossp";

create table message
(
    id             uuid default gen_random_uuid() not null
        constraint message_pk
            primary key,
    timestamp      timestamptz default now()      not null,
    message        text                           not null
);

create table message_peer_user_to_user
(
    from_user_id uuid not null,
    message_id   uuid not null
        constraint message_peer_user_to_user_pk
            primary key
        constraint message_from_user_to_user_message_fk
            references message,
    to_user_id   uuid not null
);

create table message_read_user_to_user
(
    from_user_id         uuid not null,
    to_user_id           uuid not null,
    last_read_message_id uuid not null
        constraint message_read_user_to_user_message_fk
            references message,
    constraint message_read_user_to_user_pk
        primary key (from_user_id, to_user_id)
);