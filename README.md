# Bug reproduction

Using custom marshallers on an insert input results in a new
["premarshal" function](https://github.com/Khan/genqlient/blob/v0.5.0/generate/marshal.go.tmpl#L31-L39)
that ignores JSON fields like `omitempty`.

This results in a bug where writers to a table also need write permissions on
any
[related](https://hasura.io/docs/latest/schema/postgres/table-relationships/index/)
tables.

## The Tables

We have two tables: `person` and `address`.

The `person` table has a relationship to the `address` table.

### The Roles

- The `person-writer` role has full permissions on the `person` table.
- The `address-reader` role has read permissions on the `address` table.
- The application uses both roles via an Inherited Role.

## The Bug

This is, at its core, a permissions bug.

Imagine a Postgres database where tables have
[relationships](https://hasura.io/docs/latest/schema/postgres/table-relationships/index/)
to each other, but each table requires different permissions.

In order to use an "insert input" for one table, you'll need write permissions
for both tables, even if you only intend on inserting into one table and not the
other.
