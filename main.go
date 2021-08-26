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
	// hostName := os.Getenv("DB_HOST")
	// user := os.Getenv("DB_USER")
	// password := os.Getenv("DB_PASSWORD")
	// dbname := os.Getenv("DB_NAME")
	// port := os.Getenv("PORT")

	hostName := "localhost"
	user := "postgres"
	password := "ihsan123"
	dbname := "online_store"
	port := "5432"

	// buat struct untuk keperluan credential database
	db := database.Database{
		Hostname: hostName,
		User:     user,
		Password: password,
		DBName:   dbname,
		Port:     port,
	}
	//buat instance database
	gormDB, err := database.NewPostgresDatabase(db)
	if err != nil {
		panic(err)
	}

	// buat instance fiber
	app := fiber.New()

	//setting cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: false,
		ExposeHeaders:    "*",
	}))

	//tambahkan log
	app.Use(logger.New(), recover.New())

	//Proses dependency injection

	/*
		repository merupakan layer yang melakukan query ke database
	*/
	orderRepo := repository.NewOrderRepository(gormDB)
	itemRepo := repository.NewItemRepository(gormDB)

	/*
		Service merupakan layer bussines logic yang menghubungkan handler dan repository.
	*/
	orderService := service.NewOrderService(orderRepo, itemRepo)

	/*
		Handler merupakan layer yang pertama kali menerima data (params/body) yang dikirimkan client.
		Dalam layer ini juga didefinisikan api-api yang tersedia.
	*/
	handler.NewOrderHandler(app, orderService)

	//jalankan goroutine cronjob untuk melakukan pengecekan expired order
	cron := make(chan error)

	go func(cron chan error) { cron <- orderService.CheckExpiredCheckout(1) }(cron)

	// run and listen app
	app.Listen(":8080")
}
