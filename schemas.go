package manager

import "fmt"

type Schemas []*Schema

// All is an alias for String
func (schemas Schemas) All() string {
	return schemas.String()
}

func (schemas Schemas) Get(name string) *Schema {
	for _, schema := range schemas {
		if schema.Name == name {
			return schema
		}
	}
	return nil
}

func (schemas Schemas) Has(name string) bool {
	for _, schema := range schemas {
		if schema.Name == name {
			return true
		}
	}
	return false
}

// List returns the schemas
func (schemas Schemas) List() string {
	if schemas == nil {
		return "No apps"
	}
	out := "Available schemas:\n"
	for _, schema := range schemas {
		out += fmt.Sprintf(" * %s\n", schema.Name)
	}
	return out
}

func (schemas Schemas) String() (out string) {
	for _, schema := range schemas {
		out += schema.String()
	}
	return
}

func (schemas Schemas) SQL(all bool, names ...string) string {
	if all {
		return schemas.String()
	}
	if names == nil {
		return schemas.List()
	}
	var invalid bool
	for _, name := range names {
		invalid = invalid || !schemas.Has(name)
	}
	if invalid {
		return schemas.List()
	}
	out := ""
	for _, name := range names {
		out += fmt.Sprintf("%s\n", schemas.Get(name))
	}
	return out
}

// Maintain a default list of tables
var defaults = Schemas{}

func All() string {
	return defaults.All()
}

// Defaults returns the default Schemas
func Defaults() Schemas {
	return defaults
}

func Get(name string) *Schema {
	return defaults.Get(name)
}

func Has(name string) bool {
	return defaults.Has(name)
}

// List call List on the default Schemas
func List() string {
	return defaults.List()
}

// String calls String on the default Schemas
func String() string {
	return defaults.String()
}

func SQL(all bool, names ...string) string {
	return defaults.SQL(all, names...)
}
