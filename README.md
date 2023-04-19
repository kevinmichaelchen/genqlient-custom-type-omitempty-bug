When you use custom type bindings for custom scalars (for example, for a
`tsrzrange` in Postgres), genqlient will create a
["premarshal" function](https://github.com/Khan/genqlient/blob/v0.5.0/generate/marshal.go.tmpl#L31-L39)
that ignores JSON fields like `omitempty`.
