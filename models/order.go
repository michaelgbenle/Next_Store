package models

type Status struct {
	Pending   bool `json:"pending,omitempty"`
	Canceled  bool `json:"canceled,omitempty"`
	Completed bool `json:"completed,omitempty"`
}

type Cart struct {
	Id         uint   `json:"id"gorm:"primary_key"`
	Name       string `json:"name" gorm:"name"`
	Price      int    `json:"price" gorm:"price"`
	Quantity   int    `json:"quantity"gorm:"quantity"`
	Image      string `json:"image" gorm:"image"`
	ProductID  uint   `json:"productID" gorm:"productID"`
	Buyer      Buyer  `json:"buyer"gorm:"buyer"`
	BuyerID    uint   `json:"buyerID" gorm:"buyerID"`
	TotalPrice int    `json:"totalPrice"gorm:"totalPrice"`
}


