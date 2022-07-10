CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create type FLAGTYPE as enum ('BOOLEAN', 'STRING', 'NUMBER');

create table account (
    id uuid default uuid_generate_v4() primary key,
    account_name text not null,
    account_description text not null,
    created_on timestamptz not null default now(),
    modified_on timestamptz not null default now()
);

create table token (
    id uuid default uuid_generate_v4() primary key,
    token uuid default uuid_generate_v4(),
    account_id uuid references account(id) on delete cascade,
    read_only boolean default true,
    created_on timestamptz not null default now(),
    modified_on timestamptz not null default now()
);

create table project (
    id uuid default uuid_generate_v4() primary key,
    account_id uuid references account(id) on delete cascade,
    project_name text not null,
    project_description text not null,
    created_on timestamptz not null default now(),
    modified_on timestamptz not null default now()
);

create table flag (
    id uuid default uuid_generate_v4() primary key,
    project_id uuid references project(id) on delete cascade,
    account_id uuid references account(id) on delete cascade,
    flag_key text not null,
    flag_type FLAGTYPE not null,
    flag_value text not null,
    created_on timestamptz not null default now(),
    modified_on timestamptz not null default now(),
    unique (project_id, flag_key)
);

create or replace function update_modified_on() returns trigger as $$
    begin
       NEW.modified_on := now();
       return NEW;
    end;
$$ language plpgsql;

create trigger account_modified_on
    after update or insert
    on account
    for each row
    execute procedure update_modified_on();

create trigger token_modified_on
    after update or insert
    on token
    for each row
execute procedure update_modified_on();

create trigger project_modified_on
    after update or insert
    on project
    for each row
execute procedure update_modified_on();

create trigger flag_modified_on
    after update or insert
    on flag
    for each row
execute procedure update_modified_on();

