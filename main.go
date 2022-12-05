package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

// db gorm as Database
var db *gorm.DB
var err error

// Table from product
type Product struct {
	ID    int             `json:"id"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16,2)"`
}

// Response array of product
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:bagaskaramadhan97@/go_rest_api_crud?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Cannot Connect to Database", err)
	} else {
		log.Println("Connect to Database")
	}

	// After connect table auto create
	db.AutoMigrate(&Product{})
	handleRequest()
}

func handleRequest() {
	log.Println("Start the development server at http://127.0.0.1:3000")

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)

	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}
