# Entxbulk

Adds a `BulkUpdate` method to Ent schema clients to perform multiple DB updates with a single SQL query.

## Example

Ent client interactions

```go
client.User.BulkUpdate(ctx, []ent.UserBulkUpdate{
	{
		ID:        13,
		HubspotId: api.ToPtr("abc123"),
	},
	{
		ID:        14,
		HubspotId: api.ToPtr("xyz789"),
	},
})
```

Executed SQL example

```sql
UPDATE users
SET hubspot_id = data.hubspot_id
FROM (VALUES ($1::bigint, $2), ($3::bigint, $4)) AS data(id, hubspot_id)
WHERE users.id = data.id;
```

## Configuration

In entc.go

```go
bulkUpdateExtension := entxbulk.NewExtension()
if err := entc.Generate("../ent/schema", &gen.Config{}, entc.Extensions(bulkUpdateExtension)); err != nil {
	panic(err)
}
```

For a configuration with custom data type mappings do something like

```go
entxbulk.NewExtension(
	entxbulk.WithTypeCasts(map[string]string{
		"time.Time":        "timestamptz",
		"*pgtype.Interval": "interval",
		"int":              "bigint",
	}),
	entxbulk.WithGoTypes(map[string]string{
		"interval":    "*pgtype.Interval",
		"timestamptz": "time.Time",
		"bigint":      "int",
	}),
)
```
