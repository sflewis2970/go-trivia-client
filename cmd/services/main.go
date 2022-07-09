package main

import (
	"log"
	"net/http"

	"github.com/sflewis2970/go-trivia-client/controllers"
)

func main() {
	controllers.New()

	log.Print("Web Server is ready!")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
