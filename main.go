package main

import (
	"example/todolist2/app"
	"log"
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	m := app.MakeHandler()

	// app이 종료되기전에 Close 호출
	defer m.Close()
	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("Started App")
	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
