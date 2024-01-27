package migrations

import (
	"context"
	"database/sql"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateEventsTable, downCreateEventsTable)
}

func upCreateEventsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "create table `"+storage.EventTableName+"` ( "+
		"id varchar(50) primary key,"+
		"title varchar(255),"+
		"date datetime,"+
		"date_end datetime,"+
		"description varchar(255),"+
		"owner_id int,"+
		"notification bigint"+
		");",
	)
	if err != nil {
		return err
	}

	return nil
}

func downCreateEventsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "drop table `"+storage.EventTableName+"`;")
	if err != nil {
		return err
	}

	return nil
}
