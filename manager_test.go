package manager

import (
	"os"
	"testing"

	"github.com/aodin/config"
	"github.com/aodin/fields"
	"github.com/aodin/sol"
	"github.com/aodin/sol/postgres"
	"github.com/aodin/sol/types"
)

// Test schemas
var Items = New(postgres.Table("items",
	fields.Serial{},
	sol.Column("name", types.Varchar().Limit(32).NotNull()),
	sol.Column("is_free", types.Boolean().NotNull()),
	sol.PrimaryKey("id"),
	fields.Timestamp{},
))

// Serial and Timestamp will be set by the database
type Item struct {
	fields.Serial
	Name   string
	IsFree bool
	fields.Timestamp
}

// Item should implement the Manageable interface
var _ Manageable = Item{}

func NewItem(name string) Item {
	return Item{Name: name}
}

var travisCI = config.Database{
	Driver:  "postgres",
	Host:    "localhost",
	Port:    5432,
	Name:    "travis_ci_test",
	User:    "postgres",
	SSLMode: "disable",
}

// getConfigOrUseTravis returns the parsed db.json if it exists or the
// travisCI config if it does not
func getConfigOrUseTravis() (config.Database, error) {
	conf, err := config.ParseDatabasePath("./db.json")
	if os.IsNotExist(err) {
		return travisCI, nil
	}
	return conf, err
}

func TestManager(t *testing.T) {
	conf, err := getConfigOrUseTravis()
	if err != nil {
		t.Fatalf("Failed to parse database config: %s", err)
	}

	conn, err := sol.Open(conf.Credentials())
	if err != nil {
		t.Fatalf("Failed to connect to a PostGres instance: %s", err)
	}
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Failed to begin new transaction: %s", err)
	}
	defer tx.Rollback()

	Items.Use(tx)
	tx.Query(Items.Create().Temporary().IfNotExists())

	// Create a new item
	a := NewItem("a")
	Items.Save(&a)
	if !a.Exists() {
		t.Errorf("Save failed to set item ID")
	}
	if a.CreatedAt.IsZero() {
		t.Errorf("Save failed to set created_at timestamp")
	}

	// TODO Update via Save
	// UpdateValues
	a.Name = "b"
	if err = Items.UpdateValues(a, sol.Values{"name": a.Name}); err != nil {
		t.Errorf("UpdateValues should not error: %s", err)
	}

	var b Item
	Items.Get(&b, a.ID)
	if a.ID != b.ID {
		t.Errorf("b should have the ID as a")
	}
	if b.Name != "b" {
		t.Errorf("b's name should be 'b'")
	}
}

func TestManager_Filter(t *testing.T) {
	conf, err := getConfigOrUseTravis()
	if err != nil {
		t.Fatalf("Failed to parse database config: %s", err)
	}

	conn, err := sol.Open(conf.Credentials())
	if err != nil {
		t.Fatalf("Failed to connect to a PostGres instance: %s", err)
	}
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Failed to begin new transaction: %s", err)
	}
	defer tx.Rollback()

	// Filter by only free items
	FreeItems := Items.Filter(Items.C("is_free").Equals(true))

	FreeItems.Use(tx)
	tx.Query(Items.Create().Temporary().IfNotExists())

	// Create free and non-free items
	items := []Item{
		{Name: "A"},
		{Name: "B", IsFree: true},
	}

	if err = FreeItems.BulkCreate(&items); err != nil {
		t.Fatalf("BulkCreate should not error: %s", err)
	}

	// TODO All items should have their ids set

	var free []Item
	if err = FreeItems.All(&free); err != nil {
		t.Fatalf("All should not error: %s", err)
	}
	if len(free) != 1 {
		t.Fatalf("Unexpected length of free items: %d != 1", len(free))
	}
}
