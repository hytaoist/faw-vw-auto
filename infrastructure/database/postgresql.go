package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/hytaoist/autosignin/internal/log"
	"github.com/pkg/errors"
)

type Psql struct {
	db *sql.DB
}

func NewPsql() *Psql {
	db, err := sql.Open("sqlite3", "FAWVW.db")
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	return &Psql{db}
}

func (p *Psql) Versions(product string) ([]string, error) {
	query := `
		  SELECT DISTINCT j.version
		    FROM job AS j
		   WHERE j.product = $1
		ORDER BY j.version
	`
	rows, err := p.db.Query(query, product)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	versions := ([]string)(nil)
	v := ""
	for rows.Next() {
		err = rows.Scan(&v)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		versions = append(versions, v)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return versions, nil
}
