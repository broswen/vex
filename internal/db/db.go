package db

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"log"
	"strings"
)

type Database struct {
	*pgx.Conn
}

func InitDB(ctx context.Context, dsn string) (*Database, error) {
	db, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidData  = errors.New("invalid data")
	ErrKeyNotUnique = errors.New("flag key not unique")
	ErrUnknown      = errors.New("unknown error")
)

func PgError(err error) error {
	if err == nil {
		return nil
	}
	log.Println(err.Error())

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		log.Println(pgError.Code)
		log.Println(pgError.Detail)
		log.Println(pgError.Message)
		log.Println(pgError.ConstraintName)
		log.Println(pgError.TableName)
		log.Println(pgError.ColumnName)

		//https://www.postgresql.org/docs/11/errcodes-appendix.html
		//convert postgres error codes to user friendly errors
		switch {
		case strings.HasPrefix(pgError.Code, "22"):
			return ErrInvalidData
		case strings.HasPrefix(pgError.Code, "23"):
			return ErrKeyNotUnique
		default:
			return ErrUnknown
		}
	}

	//pgx returns ErrNoRows
	if err == pgx.ErrNoRows {
		return ErrNotFound
	}

	return err
}
