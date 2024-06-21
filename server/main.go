package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	db "github.com/siiickok/json-api/db"
	handlers "github.com/siiickok/json-api/handlers"
)

func main() {
	godotenv.Load()
	port := ":" + os.Getenv("PORT")
	if port == "" {
		port = "5050"
	}

	dbPort := ":" + os.Getenv("DBPORT")
	if dbPort == "" {
		port = "5432"
	}

    connStr := fmt.Sprintf("user=postgres dbname=test port=%v password=sid sslmode=disable", dbPort)

	log.Println("[INFO]", "connecting to db at", dbPort)
	d, err := db.Connect(connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Disconnect(d)
	log.Println("[INFO]", "connected to the database")

	mux := http.NewServeMux()

	sv := handlers.NewServer(port, d, mux)
	err = sv.DB.Init()
	if err != nil {
		log.Fatalln("[ERROR]", "could not initiate database", err)
	}

	handlers.Register(&sv)

	log.Println("[INFO]", "starting server at port", sv.Addr)
	if err := http.ListenAndServe(sv.Addr, sv.Mux); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
