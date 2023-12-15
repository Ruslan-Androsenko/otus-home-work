package main

import "flag"

const (
	CmdVersion       = "version"
	CmdMigrationUp   = "migration-up"
	CmdMigrationDown = "migration-down"
)

// Проверяем запрошено ли отображение данных о версии сервиса.
func hasVersionCommand() bool {
	return flag.Arg(0) == CmdVersion
}

// Проверяем запрошено ли применение миграций.
func hasMigrationUpCommand() bool {
	return flag.Arg(0) == CmdMigrationUp
}

// Проверяем запрошено ли откат миграций.
func hasMigrationDownCommand() bool {
	return flag.Arg(0) == CmdMigrationDown
}
