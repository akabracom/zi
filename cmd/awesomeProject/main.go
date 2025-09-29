package main

import (
	"deposit-calculator/internal/api"
	"log"
)

func main() {
	log.Println("Deposit Calculator Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}
