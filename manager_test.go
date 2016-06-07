package manager

import (
	"os"
	"sync"
	"testing"

	"github.com/aodin/errors"
	"github.com/aodin/fields"
	"github.com/aodin/sol"
	"github.com/aodin/sol/postgres"
	"github.com/aodin/sol/types"
)

const travisCI = "host=localhost port=5432 dbname=manager_test user=postgres sslmode=disable"

var testconn *sol.DB
var once sync.Once

// getConn returns a postgres connection pool
func getConn(t *testing.T) *sol.DB {
	// Check if an ENV VAR has been set, otherwise, use travis
	credentials := os.Getenv("MANAGER_TEST")
	if credentials == "" {
		credentials = travisCI
	}

	once.Do(func() {
		var err error
		if testconn, err = sol.Open("postgres", credentials); err != nil {
			t.Fatalf("Failed to open connection: %s", err)
		}
		testconn.SetMaxOpenConns(25)
	})
	return testconn
}

func InitSchema(conn sol.Conn, tables ...sol.Tabular) {
	// Create the given schemas
	for _, table := range tables {
		if table == nil || table.Table() == nil {
			continue
		}
		conn.Query(table.Table().Create().IfNotExists().Temporary())
	}
}

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

func (item Item) Error(conn sol.Conn) *errors.Error {
	return nil
}

func (item *Item) Save(conn sol.Conn) error {
	return conn.Query(Items.Insert().Values(item).Returning(), item)
}

// Item should implement the Manageable interface
var _ Managed = &Item{}

func NewItem(name string) Item {
	return Item{Name: name}
}

func TestManager(t *testing.T) {
	tx, err := getConn(t).Begin()
	if err != nil {
		t.Fatalf("Failed to begin new transaction: %s", err)
	}
	defer tx.Rollback()

	Items.SetConn(tx)

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
	if err = Items.UpdateValues(&a, sol.Values{"name": a.Name}); err != nil {
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
	tx, err := getConn(t).Begin()
	if err != nil {
		t.Fatalf("Failed to begin new transaction: %s", err)
	}
	defer tx.Rollback()

	// Filter by only free items
	FreeItems := Items.Filter(Items.C("is_free").Equals(true))

	FreeItems.SetConn(tx)
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
