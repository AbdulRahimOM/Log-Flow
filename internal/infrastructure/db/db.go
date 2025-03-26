package db

import (
	"fmt"
	"log"
	"log-flow/internal/infrastructure/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Couldn't connect to DB. Error: %v", err)
	}
	return db

}

func connectToDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.Env.Postgres.User,
		config.Env.Postgres.Password,
		config.Env.Postgres.Host,
		config.Env.Postgres.Port,
		config.Env.Postgres.DbName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Couldn't connect to DB. Error:", err)
		return nil, err
	}
	return db, nil
}
