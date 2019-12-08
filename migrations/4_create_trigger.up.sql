create trigger trigger_put_object
before update or insert
on fda_dbs.user_objs
for each row
execute function fda_dbs.trigger_fn_put_object();
