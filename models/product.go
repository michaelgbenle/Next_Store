package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Name                 string `json:"name" gorm:"name"`
	Price                int    `json:"price" gorm:"price"`
	Quantity             int    `json:"quantity" gorm:"quantity"`
	Productcategory      string `json:"product_category" gorm:"productcategory"`
	Productimg           string `json:"productimg" gorm:"productimg"`
	TotalProductLaunched int    `json:"totalProductLaunched"gorm:"totalProductLaunched"`
}
