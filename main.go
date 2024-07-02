package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/devmukhtarr/accesscodeinv/app"
	"github.com/devmukhtarr/accesscodeinv/database"
	"github.com/joho/godotenv"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Welcome to check credit score!\n")
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	// connect to db
	database.ConnectDB()
	http.HandleFunc("/", getRoot)

	app.App()

	err = http.ListenAndServe(os.Getenv("PORT"), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
