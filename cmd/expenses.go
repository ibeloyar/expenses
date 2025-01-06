package main

import (
	"github.com/ibeloyar/expenses/internal/app"
	"github.com/ibeloyar/expenses/internal/config"
)

//	@title			Swagger Expenses API
//	@version		1.0
//	@description	This is server Expenses application.

//	@contact.name	API Support
//	@contact.email	example@mail.com

//	@license.name	MIT
//	@license.url	https://opensource.org/license/mit

// @host		0.0.0.0:7070
// @BasePath	/
// @schemes 	http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.MustLoad()

	app.Run(&cfg)
}
