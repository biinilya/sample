package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Record_20171020_001049 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Record_20171020_001049{}
	m.Created = "20171020_001049"

	migration.Register("Record_20171020_001049", m)
}

// Run the migrations
func (m *Record_20171020_001049) Up() {
	var create_stmt = `
	CREATE TABLE IF NOT EXISTS "record" (
		"id" serial NOT NULL PRIMARY KEY,
		"date" timestamp with time zone NOT NULL,
		"distance" numeric(6, 3) NOT NULL DEFAULT 0 ,
		"duration" numeric(12, 3) NOT NULL DEFAULT 0 ,
		"latitude" double precision NOT NULL DEFAULT 0 ,
		"longitude" double precision NOT NULL DEFAULT 0 ,
		"user_id" bigint references "user" ("id") NOT NULL
	);
	CREATE INDEX "record_user_id" ON "record" ("user_id");
	`
	m.SQL(create_stmt)

}

// Reverse the migrations
func (m *Record_20171020_001049) Down() {
	m.SQL(`DROP TABLE "record"`)
}
