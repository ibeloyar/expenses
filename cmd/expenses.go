package main

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/B-Dmitriy/expenses/internal/logger"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID              int
	Username        string
	Password        string
	Email           string
	Email_confirmed bool
}

func main() {
	config := config.MustLoad()

	logger := logger.SetupLogger(config.ENV)
	logger.Info("logger initialized", slog.String("env", config.ENV))
	db, err := sql.Open(config.Storage.DBDriver, config.Storage.DBUser+":"+config.Storage.DBPass+"@/"+config.Storage.DBName)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM users;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := make([]User, 0, 0)

	for rows.Next() {
		user := new(User)
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Email_confirmed); err != nil {
			panic(err)
		}
		users = append(users, *user)
	}

	for _, v := range users {
		fmt.Printf("User: %v\n", v)
	}
	// TODO: implement logger
	// TODO: implement storage
	// TODO: implement server
	// TODO: Run app
}
