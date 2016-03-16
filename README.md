Manager
=======

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/aodin/manager) [![Build Status](https://travis-ci.org/aodin/manager.svg?branch=master)](https://travis-ci.org/aodin/manager)

Totally not an ORM for Go.

For use with [sol/postgres](https://github.com/aodin/sol) and [fields](https://github.com/aodin/fields).

Adds convenience methods such as `Get`, auto-joins, and injected conditional clauses.

```go
// Serial and Timestamp will be set by the database
type Item struct {
    fields.Serial
    Name   string
    IsFree bool
    fields.Timestamp
}

// Create a Table and immediately wrap it in a manager. Since the TableElem
// is embedded all its methods are available
Items := manager.New(postgres.Table("items",
    fields.Serial{},
    sol.Column("name", types.Varchar().Limit(32).NotNull()),
    sol.Column("is_free", types.Boolean().NotNull()),
    sol.PrimaryKey("id"),
    fields.Timestamp{},
))

Items.Use(conn) // A connection must be set

item := NewItem("a")
Items.Save(&item) // Pass a pointer for database-level field assignments
log.Println(item.Exists())

Items.Get(&item, item.ID)
```

Extend the manager by embedding it in your own struct.


Happy Hacking!

aodin, 2016