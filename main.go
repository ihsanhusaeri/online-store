package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/online-store/entity"
	"github.com/online-store/handler"
	"github.com/online-store/repository"
	"github.com/online-store/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold: time.Second,     // Slow SQL threshold
			LogLevel:      gormLogger.Info, // Log level
			Colorful:      false,           // Disable color
		},
	)
	dsn := "host=localhost user=postgres password=ihsan123 dbname=online_store port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic(err)
	}
	db.AutoMigrate(
		entity.Item{},
		entity.Order{},
		entity.OrderItem{},
	)
	log.Println("Database successfully migrated")

	err = db.Exec("INSERT INTO items (id, name, price, stock, created_at, updated_at) SELECT 1, 'buku', 15000, 5, now(), now() WHERE NOT EXISTS (SELECT id FROM items WHERE \"name\" = 'buku');").Error

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

	orderRepo := repository.NewOrderRepository(db)
	itemRepo := repository.NewItemRepository(db)
	orderService := service.NewOrderService(orderRepo, itemRepo)
	handler.NewOrderHandler(app, orderService)

	// run and listen app
	app.Listen(":8080")
}
