package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"holidays/handlers"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/holidays/{year}", handlers.GetHolidaysByYear).Methods("GET")
	r.HandleFunc("/holidays/{festival}/{year}", handlers.GetHolidayByNameAndYear).Methods("GET")
	r.HandleFunc("/work/{year}", handlers.GetWorkdaysByYear).Methods("GET")

	http.Handle("/", r)
	fmt.Println("服务器启动于 :8081")
	http.ListenAndServe(":8081", nil)
}
