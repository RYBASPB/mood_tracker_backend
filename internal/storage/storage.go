package storage

import "time"

type MoodScore struct {
	Id     int64     `db:"id"`
	Score  int8      `db:"score"`
	Date   time.Time `db:"date"`
	UserId int64     `db:"user_id"`
}

type AddMoodScoreDto struct {
	Score  int64
	UserId int64
}

type User struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

type CreateUserDto struct {
	Name string
}
