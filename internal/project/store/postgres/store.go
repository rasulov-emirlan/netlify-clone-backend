package postgres

import (
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepo(conn *pgxpool.Pool) (*Repository, error) {
	if conn == nil {
		return nil, errors.New("project: connection to database can't be nil")
	}
	return &Repository{
		conn: conn,
	}, nil
}
