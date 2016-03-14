package manager

type Manageable interface {
	Exists() bool
	Keys() []interface{} // Primary keys
}
