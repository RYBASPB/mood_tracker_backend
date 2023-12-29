package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

func ConnectToDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database postgresql")
	}
	defer conn.Close(context.Background())

	return conn
}
