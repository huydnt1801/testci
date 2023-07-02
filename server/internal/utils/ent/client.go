package ent

import (
	"database/sql"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/huydnt1801/chuyende/internal/ent"
)

func NewClientFromDB(db *sql.DB) *ent.Client {
	drv := entsql.OpenDB("mysql", db)
	client := ent.NewClient(ent.Driver(drv))
	return client
}
