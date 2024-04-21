package main

import (
	"database/sql"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID              int
	Username        string
	Password        string
	Email           string
	Email_confirmed bool
}

type Config struct {
	DBDriver string `yaml:"DB_DRIVER"`
	DBName   string `yaml:"DB_NAME"`
	DBUser   string `yaml:"DB_USER"`
	DBPass   string `yaml:"DB_PASS"`
}

func readConfig() Config {
	config := Config{}

	file, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func main() {
	config := readConfig()
	db, err := sql.Open(config.DBDriver, config.DBUser+":"+config.DBPass+"@/"+config.DBName)
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
}
