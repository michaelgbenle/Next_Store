package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID            uint   `json:"ID" gorm:"autoincrement"`
	Name          string `json:"name" gorm:"name"`
	Email         string `json:"email" gorm:"email"`
	Username      string `json:"Username" gorm:"Username"`
	Password      string `json:"password,omitempty" gorm:"-"`
	Phonenumber   string `json:"phonenumber" gorm:"phonenumber"`
	PasswordHash  string `json:"-" gorm:"password-hash"`
	Address       string `json:"address" gorm:"address"`
	TimeCreated   string `json:"timeCreated" gorm:"timeCreated"`
	BankName      string `json:"bankName" gorm:"bankName"`
	AccountName   string `json:"accountName" gorm:"accountName"`
	AccountNumber string `json:"accountNumber" gorm:"accountNumber"`
}

type Buyer struct {
	User
	UserID  uint `json:"userID"gorm:"foreignkey"`
	BuyerID uint `json:"buyerID" gorm:"primarykey, autoincrement"`
}
type Seller struct {
	User
	SellerID uint `json:"sellerID" gorm:"SellerID"`
}

func (u *User) PasswordHasher() string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	u.PasswordHash = string(hashedPassword)
	return u.PasswordHash
}
