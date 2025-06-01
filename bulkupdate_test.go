package entxbulk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	_ "github.com/lib/pq"
	"github.com/mkraft/entxbulk/ent"
	"github.com/mkraft/entxbulk/ent/enttest"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestBulkUpdateMethodGenerated(t *testing.T) {
	bulkUpdateExtension := NewExtension()
	err := entc.Generate("./ent/schema", &gen.Config{}, entc.Extensions(bulkUpdateExtension))
	require.NoError(t, err)

	// Check that file exists
	buf, err := os.ReadFile(filepath.Join("ent", "bulk_update_gen.go"))
	require.NoError(t, err)

	// It has something resembling what we expect
	require.Contains(t, string(buf), "func (c *UserClient) BulkUpdate")
}

// TestUserBulkUpdate_Postgres tests that Postgres accepts the generates SQL query.
// This test relies on TestBulkUpdateMethodGenerated having run first for code generation.
func TestUserBulkUpdate_Postgres(t *testing.T) {
	ctx := context.Background()
	client := startPostgresClient(t)
	defer client.Close()

	u1 := client.User.Create().SetHubspotID("old1").SaveX(ctx)
	u2 := client.User.Create().SetHubspotID("old2").SaveX(ctx)

	newVal1 := "abc123"
	newVal2 := "xyz789"

	err := client.User.BulkUpdate(ctx, []ent.UserBulkUpdate{
		{ID: u1.ID, HubspotId: &newVal1},
		{ID: u2.ID, HubspotId: &newVal2},
	})
	require.NoError(t, err)

	uu1 := client.User.GetX(ctx, u1.ID)
	require.Equal(t, newVal1, *uu1.HubspotID)

	uu2 := client.User.GetX(ctx, u2.ID)
	require.Equal(t, newVal2, *uu2.HubspotID)
}

func startPostgresClient(t *testing.T) *ent.Client {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_PASSWORD": "pass",
				"POSTGRES_DB":       "testdb",
				"POSTGRES_USER":     "user",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { container.Terminate(ctx) })

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("host=%s port=%s user=user password=pass dbname=testdb sslmode=disable", host, port.Port())
	client := enttest.Open(t, "postgres", dsn)
	return client
}
