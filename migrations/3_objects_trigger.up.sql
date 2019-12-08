create function fda_dbs.trigger_fn_put_object() returns trigger as
$BODY$
begin
    if NEW.key <> OLD.key then
        raise notice 'CANNOT CHANGE KEY OF AN OLD OBJECT';
    end if;

    INSERT INTO fda_dbs.user_objs_revs(id, db_id, col_id, content, rev_timestamp)
    VALUES (NEW.id, NEW.db_id, NEW.col_id, NEW.content, NEW.rev_timestamp);

    return NEW;
end;
$BODY$
language plpgsql;
