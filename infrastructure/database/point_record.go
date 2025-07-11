package database

import (
	"database/sql"
	"time"

	"github.com/hytaoist/faw-vw-auto/internal/log"
	"github.com/pkg/errors"
)

func (p *Psql) InsertPointRecord(changedScore int) (string, error) {
	localtime := time.Now().Local().Format(time.DateTime)
	query := `
		INSERT INTO point_record (changed_score, create_at)
		     VALUES ($1, $2)
		  RETURNING id
	`
	jID := ""
	err := p.db.QueryRow(query, changedScore, localtime).Scan(&jID)
	if err != nil {
		log.Info(query)
		err = errors.WithStack(err)
		return "", err
	}
	return jID, nil
}

func (p *Psql) SumScore() (int16, error) {
	query := `
		select sum(changed_score) from point_record;
	`

	var total sql.NullInt16
	err := p.db.QueryRow(query).Scan(&total)
	if err != nil {
		log.Info(err)
	}

	if total.Valid {
		return total.Int16, nil
	} else {
		return 0, nil
	}
}
