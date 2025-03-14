package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
}

func StartServer() {
	http.HandleFunc("/get-raw-data", GetRawDataHandler)
	http.HandleFunc("/spatial-query", SpatialQueryHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
