package manager

import (
	"github.com/aodin/errors"
	"github.com/aodin/sol"
)

// Managed is an implementation target for database-backed struct types
type Managed interface {
	Error(sol.Conn) *errors.Error
	Exists() bool
	Keys() []interface{} // Primary keys
	Save(sol.Conn) error
}
