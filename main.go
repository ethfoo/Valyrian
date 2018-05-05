package main

import (
	"Valyrian/router"
	"log"
	"net/http"
)

func main() {
	router.Route()
	err := http.ListenAndServe(":8848", nil)
	if err != nil {
		log.Fatalln("listen and server error: ", err)
	}
}
