package database

import (
	"fmt"
	"github.com/decadevs/next_store/models"
	_ "github.com/decadevs/next_store/models"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"time"
)

//DECLARE A VARIABLE THAT CONNECTS WITH DB


//FUNCTION TO OPEN AND MIGRATE
func OpenAndMigrateDb() (*gorm.DB, error) {

	//storing values from .env file into the variables
	var username = os.Getenv("DB_USERNAME")
	var password = os.Getenv("DB_PASSWORD")
	var dbname = os.Getenv("DB_NAME")

	//open Database and connect using user details
	DBClient, err := gorm.Open("mysql", username+":"+password+"@tcp(127.0.0.1:3306)/"+dbname+"?charset=utf8mb4&parseTime=True&loc=Local")
	//error handling
	if err != nil {
		log.Println("checking database error", err)
	}

	//DBClient = db

	//calling the automigrate function
	AutoMigrate(DBClient)

	return DBClient, nil

}

//FUNCTION FOR AUTOMIGRATION USING GORM
func AutoMigrate(db *gorm.DB) {

	err := db.AutoMigrate(&models.User{},
		&models.Buyer{},
		&models.Product{},
		&models.Status{},
		&models.Seller{},
		&models.Cart{},
	)

	log.Println("checking database error", err)
}

//FUNCTION TO CHECK THE DB & FIND USERS BY THEIR EMAIL ADDRESS
func FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	DBClient, er := OpenAndMigrateDb()
	if er != nil {
		log.Println(er)
	}
	var err = DBClient.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}



//FUNCTION TO FIND SELLER BY EMAIL
func FindSellerByEmail(email string) (*models.Seller, error) {
	seller := &models.Seller{}
	DBClient, _ := OpenAndMigrateDb()
	// querying the db for seller by email
	var err = DBClient.Where("email = ?", email).First(seller).Error
	if err != nil {
		return nil, err
	}
	return seller, nil
}

//FUNCTION TO CREATE NEW USER
func CreateNewUser(user *models.User) error {
	DBClient, _ := OpenAndMigrateDb()
	err := DBClient.Create(user).Error
	return err
}

//CREATE A IN-MEMORY SELLER DATABASE-RECORD
func SellerDB() *gorm.DB {
	db, err := OpenAndMigrateDb()
	//error handling
	if err != nil {
		log.Println("checking database error", err)
	}
	//defer db.Close()

	//Seller details
	UserOne := models.User{
		ID:            1,
		Name:          os.Getenv("SELLER_NAME"),
		Email:         os.Getenv("SELLER_EMAIL"),
		Username:      os.Getenv("SELLER_USERNAME"),
		Password:      os.Getenv("SELLER_PASSWORD"),
		Address:       os.Getenv("SELLER_ADDRESS"),
		AccountName:   os.Getenv("SELLER_ACCOUNTNAME"),
		AccountNumber: os.Getenv("SELLER_ACCOUNT_NUMBER"),
		Phonenumber:   os.Getenv("SELLER_PHONENUMBER"),
		BankName:      os.Getenv("SELLER_BANKNAME"),
		PasswordHash:  "",
		TimeCreated:   time.Now().Format("20-02-2021, 23:12"),
	}
	UserOne.PasswordHash = UserOne.PasswordHasher()

	//CREATE SELLER
	var Seller = models.Seller{
		UserOne,
		1,
	}

	result := db.Create(&Seller)
	err = result.Error
	if err != nil {
		log.Printf("Error creating seller record : %v", err)
	}
	rowsAffected := result.RowsAffected
	fmt.Println("Number of rows affected", rowsAffected)

	return db
}
