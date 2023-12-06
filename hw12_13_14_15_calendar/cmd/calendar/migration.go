package main

import (
	"context"

	"github.com/pressly/goose/v3"

	_ "github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/migrations"
)

func migrationUp(ctx context.Context) {
	err := goose.UpContext(ctx, dbConn, MigrationsDirectory)
	if err != nil {
		logg.Fatalf("Failed to apply migrations. Error: %v", err)
	}
}

func migrationDown(ctx context.Context) {
	err := goose.DownContext(ctx, dbConn, MigrationsDirectory)
	if err != nil {
		logg.Fatalf("Failed to rollback migrations. Error: %v", err)
	}
}
