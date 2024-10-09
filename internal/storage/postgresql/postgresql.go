package postgresql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mood_tracker/internal/storage"
	"os"
)

type Storage struct {
	db *pgxpool.Pool
}

func ConnectToDB() (*Storage, error) {
	const op = "storage.postgresql.ConnectToDB"

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("%s: Couldn't create pgxpool", op)
	}
	_, err = pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: Couldn't connect to database", op)
	}

	return &Storage{pool}, err
}

func (s *Storage) CloseConnection() {
	s.db.Close()
}

func (s *Storage) AddMoodScore(dto storage.AddMoodScoreDto) (moodScoreId int64, err error) {
	const op = "storage.postgresql.AddMoodScore"
	sql, args, err := sq.
		Insert("mood_scores").
		Columns("user_id", "score").
		Values(dto.UserId, dto.Score).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	queryResult := s.db.QueryRow(context.Background(), sql, args...)

	if err := queryResult.Scan(&moodScoreId); err != nil {
		return 0, fmt.Errorf("%s: scan: %w", op, err)
	}
	return
}

func (s *Storage) GetMoodScoresByUserId(userId int64) (moodScores []storage.MoodScore, err error) {
	const op = "storage.postgresql.GetMoodScoresByUserId"
	sql, args, err := sq.
		Select("id", "score", "date", "user_id").
		From("mood_scores").
		Where("user_id = ?", userId).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows, err := s.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed query: %w", op, err)
	}
	moodScores, err = pgx.CollectRows(rows, pgx.RowToStructByName[storage.MoodScore])
	if err != nil {
		return nil, fmt.Errorf("%s: failed collecting rows: %w", op, err)
	}
	return
}

func (s *Storage) AddUser(dto storage.CreateUserDto) (userId int64, err error) {
	const op = "storage.postgresql.AddUser"
	sql, args, err := sq.
		Insert("users").
		Columns("name").
		Values(dto.Name).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	queryResult := s.db.QueryRow(context.Background(), sql, args...)
	err = queryResult.Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("%s: scan: %w", op, err)
	}
	return
}
