package api

import (
	"deposit-calculator/internal/app/handler"
	"deposit-calculator/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting Deposit Calculator server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	// Маршруты
	r.GET("/", handler.GetMonths)                     // Главная страница
	r.GET("/month/:id", handler.GetMonth)             // Страница месяца
	r.GET("/calculation", handler.GetCalculationPage) // Страница расчета вклада

	r.Run()
	log.Println("Server down")
}
