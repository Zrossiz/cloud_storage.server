package db

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

func Migrate(db *sql.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "20230612120000_create_users_table",
				Up: []string{
					`CREATE TABLE users (
						id SERIAL PRIMARY KEY,
						name VARCHAR(100) NOT NULL,
						email VARCHAR(100) UNIQUE NOT NULL,
						password VARCHAR(255) NOT NULL,
						created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
						updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
						deleted_at TIMESTAMPTZ
					)`,
				},
				Down: []string{
					"DROP TABLE users",
				},
			},
		},
	}

	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}
