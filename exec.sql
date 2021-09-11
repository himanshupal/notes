create table if not exists users (id uuid primary key default uuid_generate_v4(), username text unique not null, first_name text, last_name text, password text not null, email text, created_at timestamp default now())

create table notes (id uuid primary key default uuid_generate_v4(), title text not null, tag text[], content text, author uuid, created_at timestamp default now(), constraint fk_author foreign key(author) references users(id))
