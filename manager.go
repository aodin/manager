package manager

import (
	"log"

	"github.com/aodin/sol"
	"github.com/aodin/sol/postgres"
)

// Manager wraps a table and provides helper methods
type Manager struct {
	*postgres.TableElem
	conn sol.Conn
	// Clauses to inject into various statements
	onDelete, onSelect, onUpdate []sol.Clause
}

// AddTo adds this manager to the given App
func (m Manager) AddTo(app *App) error {
	return app.Add(m)
}

// TODO default select
func (m *Manager) All(dest interface{}) error {
	return m.Query(m.Select(), dest)
}

// BulkCreate allow the creation of multiple objects
func (m *Manager) BulkCreate(obj interface{}) error {
	return m.Query(m.Insert().Values(obj).Returning(), obj)
}

// Save creates a new object, unless Exists() == true, then it updates
func (m *Manager) Save(obj Manageable) error {
	// TODO Always check if the connection has been set?
	if obj.Exists() {
		return nil
	}
	// TODO Error if passed a non-pointer? Or silently remove the Returning?
	return m.Query(m.Insert().Values(obj).Returning(), obj)
}

// GetBy gets an object by the given field and value
func (m *Manager) GetBy(dest interface{}, k string, v interface{}) error {
	return m.Query(m.Select().Where(m.C(k).Equals(v)).Limit(1), dest)
}

func (m *Manager) wherePK(keys ...interface{}) []sol.Clause {
	// TODO error?
	// given keys should match table's primary keys
	// TODO Any way to check type or guarantee order?
	pks := m.PrimaryKey()
	if len(pks) != len(keys) {
		// TODO Only panic if the connection is a panic connection
		log.Panicf(
			"table %s has %d primary keys - %d were given",
			m.TableElem.Name(), len(pks), len(keys),
		)
	}
	clauses := make([]sol.Clause, len(pks))
	for i, col := range pks {
		clauses[i] = m.C(col).Equals(keys[i])
	}
	return clauses
}

func (m *Manager) Get(dest interface{}, keys ...interface{}) error {
	clauses := m.wherePK(keys...)
	return m.Query(m.Select().Where(clauses...).Limit(1), dest)
}

// Conn returns the manager's connection
// TODO GetConn? What methods of the connection should be available?
func (m *Manager) Conn() sol.Conn {
	return m.conn
}

// Query allows clauses to be injected into various statements
func (m *Manager) Query(stmt sol.Executable, dest ...interface{}) error {
	// We only care about events that change things
	// If the statement has a conditional, apply additional selections
	// TODO switches can't handle new statement types
	switch t := stmt.(type) {
	case sol.DeleteStmt:
		for _, clause := range m.onDelete {
			t.AddConditional(clause)
		}
		stmt = t
	case sol.InsertStmt, postgres.InsertStmt:
		// Do nothing
	case sol.UpdateStmt:
		for _, clause := range m.onUpdate {
			t.AddConditional(clause)
		}
		stmt = t
	case sol.SelectStmt:
		for _, clause := range m.onSelect {
			t.AddConditional(clause)
		}
		stmt = t
	}
	return m.conn.Query(stmt, dest...)
}

// Clause builders
// ----

// FilterDelete injects a clause whenever a DELETE statement is queried
func (m Manager) FilterDelete(clauses ...sol.Clause) Manager {
	m.onDelete = append(m.onDelete, clauses...)
	return m
}

// FilterUpdate injects a clause whenever an UPDATE statement is queried
func (m Manager) FilterUpdate(clauses ...sol.Clause) Manager {
	m.onUpdate = append(m.onUpdate, clauses...)
	return m
}

// FilterSelect injects a clause whenever a SELECT statement is queried
func (m Manager) FilterSelect(clauses ...sol.Clause) Manager {
	m.onSelect = append(m.onSelect, clauses...)
	return m
}

// Filter injects a clause whenever a DELETE, UPDATE, or SELECT
// statement is queried
func (m Manager) Filter(clauses ...sol.Clause) Manager {
	return m.FilterDelete(clauses...).FilterUpdate(clauses...).FilterSelect(clauses...)
}

// SetConnection replaces the current connection
func (m *Manager) SetConn(conn sol.Conn) {
	m.conn = conn
}

// UpdateValues updates the given obj with the given values
func (m *Manager) UpdateValues(obj Manageable, values ...sol.Values) error {
	clauses := m.wherePK(obj.Keys()...)

	// Merge the values
	update := sol.Values{}
	for _, v := range values {
		update = update.Merge(v)
	}
	return m.Query(m.Update().Values(update).Where(clauses...))
}

// Using returns a new instance of the Manager with the given connection
func (m Manager) Using(conn sol.Conn) Manager {
	m.conn = conn
	return m
}

// Use is an alias of Using
func (m Manager) Use(conn sol.Conn) Manager {
	return m.Using(conn)
}

// New creates a new manager
func New(table *postgres.TableElem) Manager {
	return Manager{TableElem: table}
}
