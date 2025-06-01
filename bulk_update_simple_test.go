package entxbulk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBulkUpdateSimple(t *testing.T) {
	ctx := context.Background()
	client := startPostgresClient(t)
	defer client.Close()

	u1 := client.User.Create().SetHubspotID("old1").SaveX(ctx)
	u2 := client.User.Create().SetHubspotID("old2").SaveX(ctx)

	newVal1 := "abc123"
	newVal2 := "xyz789"

	err := BulkUpdateHubspotIDs(ctx, client, []HubspotUpdate{
		{ID: int64(u1.ID), HubspotID: newVal1},
		{ID: int64(u2.ID), HubspotID: newVal2},
	})
	require.NoError(t, err)

	uu1 := client.User.GetX(ctx, u1.ID)
	require.Equal(t, newVal1, *uu1.HubspotID)

	uu2 := client.User.GetX(ctx, u2.ID)
	require.Equal(t, newVal2, *uu2.HubspotID)
}
