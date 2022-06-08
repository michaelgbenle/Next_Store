package main

import (
	"fmt"
	"github.com/decadevs/next_store/database"
	_ "github.com/decadevs/next_store/models"
	"github.com/decadevs/next_store/routes"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	_ "gorm.io/driver/mysql"
	"log"
	_ "net/http"
	"os"
)

func main() {
	//database connection
	db, er := database.OpenAndMigrateDb()
	if er != nil {
		log.Println("Error in Launching Database")
		log.Fatal(er)
	}

	defer db.Close()
	//delay database shutdown

	//seller's in-memory data
	database.SellerDB()

	//call routes /sever
	routes.CallRoutes("port", db)

	//Load env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Could not load .env file %v", err)
		os.Exit(1)
	}
}
