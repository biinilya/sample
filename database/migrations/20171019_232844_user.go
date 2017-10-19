package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type User_20171019_232844 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &User_20171019_232844{}
	m.Created = "20171019_232844"

	migration.Register("User_20171019_232844", m)
}

// Run the migrations
func (m *User_20171019_232844) Up() {
	var create_stmt = `
	CREATE TABLE IF NOT EXISTS "user" (
		"id" serial NOT NULL PRIMARY KEY,
		"key" varchar(128) NOT NULL UNIQUE,
		"secret" varchar(128) NOT NULL DEFAULT '' ,
		"created" timestamp with time zone NOT NULL,
		"updated" timestamp with time zone NOT NULL
	);
	CREATE INDEX "user_key" ON "user" ("key");

	CREATE TABLE IF NOT EXISTS "user_permissions" (
		"id" serial NOT NULL PRIMARY KEY,
		"user_id" bigint references "user" ("id") NOT NULL,
		"permission_id" bigint references "permission" ("id") NOT NULL,
		unique ("user_id", "permission_id")
	);
	`
	m.SQL(create_stmt)
}

// Reverse the migrations
func (m *User_20171019_232844) Down() {
	m.SQL(`DROP TABLE "user_permissions"; DROP TABLE "user"`)
}
