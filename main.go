package main

import (
	"BCHat/routes"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	conn, dbErr := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if dbErr != nil {
		log.Fatal("Failed to connect to database")
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal("Closed connection to the database")
		}
	}(conn, context.Background())
}

func main() {
	routes.GetRoutes()
}
