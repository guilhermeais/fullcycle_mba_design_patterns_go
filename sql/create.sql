-- create extension if not exists "uuid-ossp";

create table contract (
    id uuid not null default uuid_generate_v4() primary key,
    description text,
    amount numeric,
    periods integer,
    date timestamp  
);

create table payment (
    id uuid not null default uuid_generate_v4() primary key,
    contract_id uuid references contract (id),
    amount numeric,
    date timestamp
);