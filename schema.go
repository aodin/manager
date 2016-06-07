package manager

import (
	"fmt"

	"github.com/aodin/sol"
)

type Schema struct {
	Name   string
	Tables []sol.Tabular
}

func (schema Schema) String() (out string) {
	if schema.Tables == nil {
		return fmt.Sprintf("-- %s (no tables)\n", schema.Name)
	} else {
		out = fmt.Sprintf("-- %s\n", schema.Name)
		for _, table := range schema.Tables {
			out += fmt.Sprintf("%s\n", table.Table().Create())
		}
	}
	return out
}

// Add adds a table to the Schema. It will error if a table with the same
// name already exists in the Schema
func (schema *Schema) Add(table sol.Tabular) error {
	if table == nil || table.Table() == nil {
		return fmt.Errorf(
			"manager: cannot add a nil table to schema %s", schema.Name,
		)
	}
	for _, existing := range schema.Tables {
		if table.Name() == existing.Name() {
			return fmt.Errorf(
				"manager: a table named %s already exists in schema %s",
				table.Name(), schema.Name,
			)
		}
	}

	schema.Tables = append(schema.Tables, table)
	return nil
}

// NewSchema creates a new schema - a simple way to aggregate tables
func NewSchema(name string) *Schema {
	schema := &Schema{Name: name}
	defaults = append(defaults, schema)
	return schema
}
