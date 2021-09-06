package main

import (
	"dewble/todos/app"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	m := app.MakeHandler("./test.db")

	// app이 종료되기전에 Close 호출
	defer m.Close()

	log.Println("Started App")
	err := http.ListenAndServe(":port", m)
	if err != nil {
		panic(err)
	}
}
