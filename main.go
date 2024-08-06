package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"holidays/handlers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/holidays/{year}", handlers.GetHolidaysByYear).Methods("GET")
	r.HandleFunc("/holidays/{festival}/{year}", handlers.GetHolidayByNameAndYear).Methods("GET")
	r.HandleFunc("/work/{year}", handlers.GetWorkdaysByYear).Methods("GET")

	http.Handle("/", r)
	fmt.Println("服务器启动于 :8080")
	http.ListenAndServe(":8080", nil)
}
