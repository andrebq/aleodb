create schema fda_users;

create table fda_users.users(
    id uuid not null,
    name text not null,
    role text not null,
    extra jsonb not null default '{}'::jsonb
);

alter table fda_users.users
    add constraint pk_user_id
    primary key (id);

create unique index unq_fda_users_users_name
    on fda_users.users(name);

create schema fda_dbs;
create table fda_dbs.user_dbs(
    id uuid not null,
    owner_id uuid not null,
    name text not null,
    extra jsonb not null default '{}'::JSONB
);
alter table fda_dbs.user_dbs
    add constraint pk_db_id
    primary key (id);

create unique index unq_fda_dbs_user_dbs_name
    on fda_dbs.user_dbs(owner_id, name);

alter table fda_dbs.user_dbs
    add constraint fk_db_id
    foreign key (owner_id)
    references fda_users.users(id);
