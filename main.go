package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/online-store/database"
	"github.com/online-store/handler"
	"github.com/online-store/repository"
	"github.com/online-store/service"
)

func main() {
	// hostName := os.Getenv("HOSTNAME")
	// user := os.Getenv("USER")
	// password := os.Getenv("PASSWORD")
	// dbname := os.Getenv("DB_NAME")
	// port := os.Getenv("PORT")
	// sslmode := os.Getenv("SSL_MODE")
	// timeZone := os.Getenv("TIME_ZONE")

	hostName := "localhost"
	user := "postgres"
	password := "ihsan123"
	dbname := "online_store"
	port := "5432"
	sslmode := "disable"
	timeZone := "Asia/Jakarta"

	db := database.Database{
		Hostname: hostName,
		User:     user,
		Password: password,
		DBName:   dbname,
		Port:     port,
		SSLMode:  sslmode,
		TimeZone: timeZone,
	}
	gormDB, err := database.NewGormDatabase(db)
	if err != nil {
		panic(err)
	}

	// create router instance and apply appropriate middlewares
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: false,
		ExposeHeaders:    "*",
	}))
	app.Use(logger.New(), recover.New())

	orderRepo := repository.NewOrderRepository(gormDB)
	itemRepo := repository.NewItemRepository(gormDB)
	orderService := service.NewOrderService(orderRepo, itemRepo)
	handler.NewOrderHandler(app, orderService)

	cron := make(chan error)

	go func(cron chan error) { cron <- orderService.CheckExpiredCheckout(1) }(cron)

	// run and listen app
	app.Listen(":8080")
}
