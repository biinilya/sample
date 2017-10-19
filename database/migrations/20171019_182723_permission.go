package main

import "github.com/astaxie/beego/migration"

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
	var create_table = `
	CREATE TABLE IF NOT EXISTS "permission" (
		"id" serial NOT NULL PRIMARY KEY,
		"title" varchar(128) NOT NULL UNIQUE,
		"description" text NOT NULL DEFAULT ''
	);
	`
	var embedded_permissions = `
	INSERT INTO
		permission (title, description)
	VALUES
		('user', 'can only access own records'),
		('manager', 'can only access users'),
		('admin', 'can access both all records and users')
	`

	m.SQL(create_table)
	m.SQL(embedded_permissions)
}

// Reverse the migrations
func (m *Permissions_20171016_182723) Down() {
	m.SQL(`DROP TABLE permission`)
}
