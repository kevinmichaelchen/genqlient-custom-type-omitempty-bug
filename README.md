When you use custom type bindings for custom scalars (for example, for a
`tsrzrange` in Postgres), genqlient will create a
["premarshal" function](https://github.com/Khan/genqlient/blob/v0.5.0/generate/marshal.go.tmpl#L31-L39)
that ignores JSON fields like `omitempty`.

### Tables

We have two tables: `person` and `address`.

The `person` table has a relationship to the `address` table.

### Roles

- The `person-writer` role has full permissions on the `person` table.
- The `address-reader` role has read permissions on the `address` table.
- The application uses both roles via an Inherited Role. 