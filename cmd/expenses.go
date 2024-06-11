package main

import (
	"github.com/B-Dmitriy/expenses/internal/app"
	"github.com/B-Dmitriy/expenses/internal/config"
)

func main() {
	cfg := config.MustLoad()

	app.Run(&cfg)
}
