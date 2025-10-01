package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DataBase struct {
	Conn *sql.DB
}

func NewDataBase(dsn string) (*DataBase, error) {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &DataBase{Conn: dbConn}, nil
}

func EnsureUser(database *DataBase, username string) (int, error) {
	var userID int
	err := database.Conn.QueryRow(
		`INSERT INTO users (username) 
		 VALUES ($1) 
		 ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username 
		 RETURNING id`,
		username,
	).Scan(&userID)
	return userID, err
}

func EnsureChannel(database *DataBase, channelName string) (int, error) {
	var channelID int
	err := database.Conn.QueryRow(
		`INSERT INTO channels (name) 
		 VALUES ($1) 
		 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name 
		 RETURNING id`,
		channelName,
	).Scan(&channelID)
	return channelID, err
}
