create table fda_dbs.user_cols(
    id uuid not null,
    db_id uuid  not null,
    name text not null
);

alter table fda_dbs.user_cols
    add constraint pk_fda_dbs_user_cols
    primary key(db_id, id);

alter table fda_dbs.user_cols
    add constraint fk_fda_dbs_user_cols_db
    foreign key (db_id)
    references fda_dbs.user_dbs(id);

create table fda_dbs.user_objs(
    id uuid not null,
    key text not null,
    db_id uuid not null,
    col_id uuid not null,
    content jsonb not null,
    rev_timestamp timestamptz not null default current_timestamp
);

create table fda_dbs.user_objs_revs(
    id uuid not null,
    db_id uuid not null,
    col_id uuid not null,
    content jsonb not null,
    rev_timestamp timestamptz not null
);

alter table fda_dbs.user_objs
    add constraint pk_fda_dbs_user_objs
    primary key (db_id, id);

alter table fda_dbs.user_objs
    add constraint fk_fda_dbs_user_objs_cols
    foreign key (col_id, db_id)
    references fda_dbs.user_cols(id, db_id);
