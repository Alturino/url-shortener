create table if not exists urls (
    id uuid primary key not null default gen_random_uuid(),
    url text not null,
    short_url varchar(5) not null,
    created_at timestamp not null default (now()),
    updated_at timestamp not null default (now()),
    visited_count int
);
