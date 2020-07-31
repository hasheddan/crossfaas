package function

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/openfaas/openfaas-cloud/sdk"
)

var (
	dbHostKey     = "endpoint"
	dbUserKey     = "username"
	dbPasswordKey = "password"

	errMessage = "Unable to service request"
)

var (
	host     string
	user     string
	password string
)

func read(key string) string {
	val, err := sdk.ReadSecret(key)
	if err != nil {
		panic(err)
	}
	return val
}

var db *gorm.DB
var err error

type ServerlessTalk struct {
	Title   string `json:"title"`
	Host    string `json:"host"`
	Viewers int    `json:"viewers"`
}

func init() {
	host = read(dbHostKey)
	user = read(dbUserKey)
	password = read(dbPasswordKey)
}

// Handle a function invocation
func Handle(w http.ResponseWriter, r *http.Request) {
	var err error
	db, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=5432 user=%s dbname=postgres password=%s", host, user, password))
	defer db.Close()

	if err != nil {
		log.Printf("Connection Failed to Open: %v", err)
	} else {
		log.Println("Connection Established")
	}

	db.AutoMigrate(&ServerlessTalk{})

	streams := []ServerlessTalk{}
	db.Find(&streams)
	fmt.Println("Request Received: /list")
	json.NewEncoder(w).Encode(streams)
}
