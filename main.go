package main

import (
	"fmt"
	"net/http"
	"profile-form/controllers"
)

func main() {
	http.HandleFunc("/", controllers.HandleForm)
	http.HandleFunc("/submit", controllers.HandleSubmit)

	fmt.Println("server running on port 3000")
	http.ListenAndServe(":3000", nil)
}
