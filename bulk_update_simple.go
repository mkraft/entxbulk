package entxbulk

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mkraft/entxbulk/ent"
	entuser "github.com/mkraft/entxbulk/ent/user"
)

const maxBatchSize = 1000

type HubspotUpdate struct {
	ID        int64
	HubspotID string
}

func BulkUpdateHubspotIDs(ctx context.Context, client *ent.Client, updates []HubspotUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	// Input validation
	for i, u := range updates {
		if u.HubspotID == "" {
			return fmt.Errorf("empty HubspotID at index %d", i)
		}
		if u.ID <= 0 {
			return fmt.Errorf("invalid ID at index %d: %d", i, u.ID)
		}
	}

	// Process in batches
	for i := 0; i < len(updates); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(updates) {
			end = len(updates)
		}

		batch := updates[i:end]

		if err := processBatch(ctx, client, batch); err != nil {
			return fmt.Errorf("error processing batch at offset %d: %w", i, err)
		}
	}

	return nil
}

func processBatch(ctx context.Context, client *ent.Client, updates []HubspotUpdate) error {
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

	_, err := client.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}
