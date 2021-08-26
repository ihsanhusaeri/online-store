package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/online-store/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Hostname string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
	TimeZone string
}

func NewGormDatabase(database Database) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // Disable color
		},
	)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", database.Hostname, database.User, database.Password, database.DBName, database.Port, database.SSLMode, database.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		return nil, err
	}
	db.AutoMigrate(
		entity.Item{},
		entity.Order{},
		entity.OrderItem{},
	)
	log.Println("Database successfully migrated")

	err = db.Exec("INSERT INTO items (id, name, price, stock, created_at, updated_at) SELECT 1, 'buku', 15000, 5, now(), now() WHERE NOT EXISTS (SELECT id FROM items WHERE \"name\" = 'buku');").Error

	if err != nil {
		return nil, err
	}
	return db, nil
}
