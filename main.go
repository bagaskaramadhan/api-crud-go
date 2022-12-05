package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
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
	ID    string          `json:"id,omitempty"`
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
	myRouter.HandleFunc("/api/v1/products", postCreateProduct).Methods("POST")
	myRouter.HandleFunc("/api/v1/products", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/v1/products/{id}", getProductDetail).Methods("GET")
	myRouter.HandleFunc("/api/v1/products/{id}", putUpdateProduct).Methods("PUT")
	myRouter.HandleFunc("/api/v1/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func postCreateProduct(w http.ResponseWriter, r *http.Request) {

	payloads, _ := ioutil.ReadAll(r.Body)
	var product Product
	product.ID = uuid.New().String()

	json.Unmarshal(payloads, &product)
	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Success to create product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	product := []Product{}

	db.Find(&product)

	res := Result{Code: 200, Data: product, Message: "Success to get products"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProductDetail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID := params["id"]
	var product Product
	db.First(&product, "id = ?", productID)

	res := Result{Code: 200, Data: product, Message: "Success to get product detail"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func putUpdateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID := params["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	// R.Body
	var productUpdate Product
	json.Unmarshal(payloads, &productUpdate)

	// Update
	var product Product
	db.Model(&product).Where("id = ?", productID).Update(productUpdate)

	res := Result{Code: 200, Data: productUpdate, Message: "Success to update product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID := params["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payloads, &product)
	productWillDelete := db.First(&product, "id = ?", productID)
	db.Delete(&product).Where("id = ?", productID)

	res := Result{Code: 200, Data: productWillDelete, Message: "Success to delete product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
