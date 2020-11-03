package main

import (
	"fmt"
	"net/http"
	"ws/initialize"
)

func main() {
	initialize.StartServer()
	http.HandleFunc("/ws", initialize.ServeHTTP)
	err := http.ListenAndServe(":9912", nil)
	if err != nil {
		fmt.Printf("http service start fail at  %s\n", err.Error())
	}
}
