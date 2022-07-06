CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create type FLAGTYPE as enum ('BOOLEAN', 'STRING', 'NUMBER');

create table account (
    id uuid default uuid_generate_v4() primary key,
    account_name text not null,
    account_description text not null
);

create table token (
    id uuid default uuid_generate_v4() primary key,
    token uuid default uuid_generate_v4(),
    account_id uuid references account(id) on delete cascade,
    read_only boolean default true
);

create table project (
    id uuid default uuid_generate_v4() primary key,
    account_id uuid references account(id) on delete cascade,
    project_name text not null,
    project_description text not null
);

create table flag (
    id uuid default uuid_generate_v4() primary key,
    project_id uuid references project(id) on delete cascade,
    account_id uuid references account(id) on delete cascade,
    flag_key text not null,
    flag_type FLAGTYPE not null,
    flag_value text not null,
    unique (project_id, flag_key)
);