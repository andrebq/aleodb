# aleodb
A simple object database (work-in-progress) to allow storage of ale objects into postgresql jsonb fields.

# Basic operation

Just call `aleodb.Register(manager)` to register a qualifed namespace `odb` into your `ale` environment.

The system assumes that `users` can have many `databases` with many `collections`.

Internally everything is reduced to a flatten structure managed by `UUID` starting at the root UUID `453632b9-dcc0-4823-815b-4f36e6ad6353`.

That means each UUID generated is unique as long as the human friendly id is unique:

- UserID = UUIDv5(rootUUID, "your unique user name or id")
- DatabaseID = UUIDv5(UserID, "name of a database owned by the given user")
- CollectionID = UUIDv5(DatabaseID, "name of a collection inside the database")
- ObjectID = UUIDv5(CollectionID, "your object key")

Thus, all objects from all users, could live in a flat space of UUIDs. SHA1 collisions are negligible in this scenario as long as you restrict the length of the human friendly text.

To allow for easy partitioning of data internally `DatabaseID` is part of all rows so you could host different logical databases on different tables and benefit from solutions like CitusDB to shard those tables into different machines.

To have a better picture of how the data is stored inside postgresql check the [migrations](./migrations) folder.

## Why the Go public API is so minimal?

The system is designed to work within Ale, reading data directly from Go could lead to some hard to fix issues.

In the future there might be a Go public API, but right now, I don't have a need for it so I will refrain from adding it.

If you want such API, I'm open to hear proposals.

## Why postgresql?

Because it is the database I'm using for the project which will use `aleodb` but any database capable of storing JSON documents (and querying them) is fine.

SQLite might be useful but I don't need to support it right now (feel free to open a PR to discuss support for multiple backends).
