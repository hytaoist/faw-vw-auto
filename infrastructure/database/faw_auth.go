package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/hytaoist/faw-vw-auto/domain"
	"github.com/hytaoist/faw-vw-auto/internal/log"
	"github.com/pkg/errors"
)

func (p *Psql) CreateFAW_Auth(ctx context.Context, auth *domain.FAWAuth) error {
	_, err := p.InsertAuth(ctx, auth.AccessToken, auth.TokenType, auth.ExpiresIn)
	if err != nil {
		return err
	}
	return nil
}

func (p *Psql) InsertAuth(ctx context.Context, accessToken string, tokenType string, expiresIn string) (string, error) {
	localtime := time.Now().Local().Format(time.DateTime)
	query := `
		INSERT INTO faw_auth (access_token, token_type, expires_in, create_at)
		     VALUES ($1, $2, $3, $4)
		  RETURNING id
	`
	jID := ""
	err := p.db.QueryRow(query, accessToken, tokenType, expiresIn, localtime).Scan(&jID)
	if err != nil {
		log.Info(query)
		err = errors.WithStack(err)
		return "", err
	}
	return jID, nil
}

func (p *Psql) FindLatestOne() (*domain.FAWAuth, error) {
	query := `
		select access_token, token_type, expires_in from faw_auth order by id desc limit 1;
	`
	var auth domain.FAWAuth
	row := p.db.QueryRow(query)
	//Scan
	err := row.Scan(&auth.AccessToken, &auth.TokenType, &auth.ExpiresIn)
	if err != nil {
		// 处理没有找到行的情况
		if err == sql.ErrNoRows {
			return &domain.FAWAuth{}, nil
		}

		// 错误信息
		log.Info(query)
		err = errors.WithStack(err)
		return nil, err
	}
	return &auth, nil
}
