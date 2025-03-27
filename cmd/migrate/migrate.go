package main

import (
	"fmt"
	"log-flow/internal/domain/models"
	"log-flow/internal/infrastructure/db"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Hello, World!")

	db := db.GetDB()

	err := migrateTables(db, []models.DbTablesWithName{
		models.Job{},
		models.LogReport{},
		models.TrackedKeywordsCount{},
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Logs table migrated successfully")
}

func migrateTables(db *gorm.DB, tables []models.DbTablesWithName) error {
	for _, table := range tables {
		err := db.AutoMigrate(table)
		if err != nil {
			return fmt.Errorf("Error migrating table %s: %v", table.TableName(), err)
		}
	}

	return nil
}
