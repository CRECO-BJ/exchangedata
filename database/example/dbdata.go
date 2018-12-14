package main

import (
	//	"github.com/exchangedata/database"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Address struct {
	Address1 string `gorm:"primary_key"`
	UserID   int
}

type Email struct {
	UserRefer int
	Email     string `gorm:"primary_key"`
}

type Language struct {
	Name  string `gorm:"primary_key"`
	Users []User `gorm:"many2many:user_languages;"`
}

type User struct {
	gorm.Model
	Test
	Name            string
	BillingAddress  Address    `gorm:"foreignkey:UserID"`
	ShippingAddress Address    `gorm:"foreignkey:UserID"`
	Emails          []Email    `gorm:"foreignkey:UserRefer"`
	Languages       []Language `gorm:"many2many:user_languages;"`
}

type Test struct {
	Long  int
	Short int
}

func main() {
	name := "exdata"
	passwd := "office98"
	user := "root"

	db, err := gorm.Open("mysql", user+":"+passwd+"@tcp(127.0.0.1:3306)/"+name+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf("open database failed, %s", err)
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &Address{}, &Email{}, &Language{})

	tuser := User{
		Name:            "jinzhu",
		BillingAddress:  Address{Address1: "Billing Address - Address 1"},
		ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
		Emails: []Email{
			{Email: "jinzhu@example.com"},
			{Email: "jinzhu-2@example@example.com"},
		},
		Languages: []Language{
			Language{Name: "ZH"},
			Language{Name: "EN"},
		},
		Test: Test{2, 1000},
	}

	db.Save(&tuser)
}

// Clean the test Database:
// drop table addresses,emails,user_languages,users,languages;
