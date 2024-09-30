package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

type Storage struct {
	db *pgxpool.Pool
}

func ConnectToDB() (*Storage, error) {
	const op = "storage.postgresql.ConnectToDB"

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	_, err = pool.Exec(context.Background(), `
		create table if not exists mood_tracker(
		    id integer primary key,
		    score integer not null);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{pool}, err
}
