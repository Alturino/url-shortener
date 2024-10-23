create table if not exists urls (
    id uuid primary key not null default (gen_random_uuid()),
    url text not null,
    short_url varchar(7) unique not null default (''),
    created_at timestamp not null default (now()),
    updated_at timestamp not null default (now()),
    visited_count int not null default (0)
);

create index if not exists idx_short_url on urls (short_url);
