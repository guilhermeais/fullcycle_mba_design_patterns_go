drop schema if exists invoices_service cascade;

create schema invoices_service;

create extension if not exists "uuid-ossp";

create table invoices_service.contract (
    id uuid not null default uuid_generate_v4() primary key,
    description text,
    amount numeric,
    periods integer,
    date timestamp  
);

create table invoices_service.payment (
    id uuid not null default uuid_generate_v4() primary key,
    contract_id uuid references invoices_service.contract (id),
    amount numeric,
    date timestamp
);

insert into invoices_service.contract values ('fac05a57-7d61-4283-ab32-7696902eac44', 'prestação de serviços escolares', 6000, 12, '2024-12-19t10:00:00');
insert into invoices_service.payment values ('6355b223-fce0-4f7c-998a-1f027281e308', 'fac05a57-7d61-4283-ab32-7696902eac44', 6000, '2024-12-18t10:00:00');