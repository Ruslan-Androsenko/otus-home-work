//nolint:revive
package config

import (
	"context"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/logger"
	_ "github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/migrations"
	"github.com/pressly/goose/v3"
)

func MigrationUp(ctx context.Context, logg *logger.Logger) {
	err := goose.UpContext(ctx, dbConn, MigrationsDirectory)
	if err != nil {
		logg.Fatalf("Failed to apply migrations. Error: %v", err)
	}
}

func MigrationDown(ctx context.Context, logg *logger.Logger) {
	err := goose.DownContext(ctx, dbConn, MigrationsDirectory)
	if err != nil {
		logg.Fatalf("Failed to rollback migrations. Error: %v", err)
	}
}
