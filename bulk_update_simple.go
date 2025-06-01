package entxbulk

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mkraft/entxbulk/ent"
	entuser "github.com/mkraft/entxbulk/ent/user"
)

type HubspotUpdate struct {
	ID        int64
	HubspotID string
}

func BulkUpdateHubspotIDs(ctx context.Context, client *ent.Client, updates []HubspotUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	var (
		params     []interface{}
		valueRows  []string
		paramIndex = 1
	)

	for _, u := range updates {
		valueRows = append(valueRows, fmt.Sprintf("($%d::bigint, $%d)", paramIndex, paramIndex+1))
		params = append(params, u.ID, u.HubspotID)
		paramIndex += 2
	}

	valuesSQL := fmt.Sprintf(
		"(VALUES %s) AS data(%s, %s)",
		strings.Join(valueRows, ", "),
		entuser.FieldID,
		entuser.FieldHubspotID,
	)

	query := fmt.Sprintf(`
		UPDATE %s
		SET %s = data.%s
		FROM %s
		WHERE %s.%s = data.%s`,
		entuser.Table,
		entuser.FieldHubspotID,
		entuser.FieldHubspotID,
		valuesSQL,
		entuser.Table,
		entuser.FieldID,
		entuser.FieldID,
	)

	if _, err := client.ExecContext(ctx, query, params...); err != nil {
		return err
	}

	return nil
}
