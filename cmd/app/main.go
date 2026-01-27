package main

import (
	"fmt"

	"github.com/CosmeticsShiraz/Backend/bootstrap"
	"github.com/CosmeticsShiraz/Backend/internal/domain/entity"
	"github.com/CosmeticsShiraz/Backend/internal/infrastructure/websocket"
	"github.com/CosmeticsShiraz/Backend/internal/presentation/routes"
	"github.com/CosmeticsShiraz/Backend/wire"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	ginEngine := gin.New()

	config := bootstrap.Run()

	hub := websocket.NewHub()
	go hub.Run()

	app, err := wire.InitializeApplication(config, hub)
	if err != nil {
		panic(err)
	}

	app.Database.DB.GetDB().AutoMigrate(
		&entity.Address{},
		&entity.City{},
		&entity.NotificationSetting{},
		&entity.NotificationType{},
		&entity.Notification{},
		&entity.Permission{},
		&entity.Province{},
		&entity.Role{},
		&entity.User{},
	)

	app.Seeds.AddressSeeder.SeedProvincesAndCities()
	app.Seeds.NotificationTypeSeeder.SeedNotificationTypes()
	app.Seeds.RoleSeeder.SeedRoles()
	app.Seeds.ContactType.SeedContactTypes()

	if err := app.Consumers.Register.Start(); err != nil {
		panic(err)
	}
	if err := app.Consumers.Push.Start(); err != nil {
		panic(err)
	}
	if err := app.Consumers.Email.Start(); err != nil {
		panic(err)
	}
	if err := app.Consumers.Notification.Start(); err != nil {
		panic(err)
	}

	routes.Run(ginEngine, app)

	ginEngine.Run(fmt.Sprintf(":%v", config.Env.Server.Port))
}
