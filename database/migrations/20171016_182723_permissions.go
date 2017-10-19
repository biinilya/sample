package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Permissions_20171016_182723 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Permissions_20171016_182723{}
	m.Created = "20171016_182723"

	migration.Register("Permissions_20171016_182723", m)
}

// Run the migrations
func (m *Permissions_20171016_182723) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE permissions(id serial primary key,title TEXT NOT NULL,description TEXT NOT NULL)")
}

// Reverse the migrations
func (m *Permissions_20171016_182723) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE permissions")
}
