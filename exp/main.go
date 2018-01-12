package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_"github.com/lib/pq"
	"bufio"
	"os"
	"strings"
)

type Order struct {
	gorm.Model
	UserID uint
	Amount int
	Description string
}

type User struct {
	gorm.Model
	Name string
	Email string `gorm:"not null;unique_index"`
	Orders []Order
}

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "zxcvbnm"
	dbname = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&User{}, &Order{})
	var user User
	db.Preload("Orders").First(&user)
	if db.Error != nil {
		panic(db.Error)
	}
	fmt.Println("Email ",user.Email)
	fmt.Println("Orders", user.Orders)
}

func getInfo() (name, email string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("name?")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println("Email?")
	email, _ = reader.ReadString('\n')
	email = strings.TrimSpace(email)
	return name, email
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	db.Create(&Order{
		UserID: user.ID,
		Amount: amount,
		Description: desc,
	})
	if db.Error != nil {
		panic(db.Error)
	}
}