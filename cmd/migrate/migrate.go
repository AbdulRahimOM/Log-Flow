package main

import (
	"fmt"
	"log-flow/internal/domain/models"
	"log-flow/internal/infrastructure/db"
)

func main() {
	fmt.Println("Hello, World!")

	db := db.GetDB()

	err := db.AutoMigrate(&models.LogReport{})
	if err != nil {
		fmt.Println("Error migrating logs table: ", err)
	}

	fmt.Println("Logs table migrated successfully")

}
