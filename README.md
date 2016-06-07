Manager
=======

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/aodin/manager)
[![Build Status](https://travis-ci.org/aodin/manager.svg?branch=master)](https://travis-ci.org/aodin/manager)

Totally not an ORM for Go.

For use with the SQL toolkits [Sol](https://github.com/aodin/sol) and [Fields](https://github.com/aodin/fields).

Manager adds convenience methods such as `Get`, auto-joins, and injected conditional clauses to SQL tables.

```go
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

// Create a Table and immediately wrap it in a manager. Since the TableElem
// is embedded all its methods are available
var Items = postgres.Table("items",
    fields.Serial{},
    sol.Column("name", types.Varchar().Limit(32).NotNull()),
    sol.Column("is_free", types.Boolean().NotNull()),
    sol.PrimaryKey("id"),
    fields.Timestamp{},
)

var ItemsManager = manager.New(Items)

func main() {
    Items.Use(conn) // A connection must be set

    item := NewItem("a")
    Items.Save(&item) // Pass a pointer for database-level field assignments
    log.Println(item.Exists())

    Items.Get(&item, item.ID)
}
```

Extend the manager by embedding it in your own struct.

Happy Hacking!

aodin, 2016
