package main

import (
	"fmt"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/presentation/routes"
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	ginEngine := gin.New()

	config := bootstrap.Run()

	app, err := wire.InitializeApplication(config)
	if err != nil {
		panic(err)
	}

	app.Database.DB.GetDB().AutoMigrate(
		&entity.Address{},
		&entity.City{},
		&entity.Province{},
		&entity.User{},
		&entity.Category{},
		&entity.Brand{},
		&entity.Product{},
		&entity.Currency{},
		&entity.Wallet{},
		&entity.Cart{},
		&entity.CartItem{},
		&entity.Transaction{},
	)
	
	app.Seeds.AddressSeeder.SeedProvincesAndCities()
	app.Seeds.AdminSeeder.SeedAdmins()
	app.Seeds.CurrencySeeder.SeedCurrencies()

	routes.Run(ginEngine, app)

	ginEngine.Run(fmt.Sprintf(":%v", config.Env.Server.Port))
}
