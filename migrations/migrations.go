package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrations interface {
	Up() (err error)
}

type migrations struct {
	migrate *migrate.Migrate
}

func NewMigrations(dsn, path string) (Migrations, error) {
	m, err := migrate.New(path, dsn)
	if err != nil {
		return nil, err
	}
	return &migrations{
		migrate: m,
	}, nil
}

func (m *migrations) Up() (err error) {
	err = m.migrate.Up()
	if err != nil && err.Error() == "no change" {
		err = nil
	}
	return
}
