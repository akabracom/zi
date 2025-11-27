package api

import (
    "log"

    "deposit-calculator/internal/app/handler"
    "deposit-calculator/internal/app/repository"

    "github.com/gin-gonic/gin"
)

func Run() {
    repo, err := repository.NewRepository()
    if err != nil {
        log.Fatalf("failed to init repository: %v", err)
    }

    h := handler.NewHandler(repo)
    // Предзаполнение корзины статическими данными (ID 1 и 2)
    h.PrefillCart()

    r := gin.Default()
    r.LoadHTMLGlob("templates/*")
    r.Static("/static", "./static")

    // Маршруты
    r.GET("/", h.GetMonths)
    r.GET("/month/:id", h.GetMonthDetail) // страница деталей

    r.POST("/set-sum", h.SetCommonAmount)
    r.POST("/add-to-cart", h.AddToCart)
    r.POST("/remove-from-cart/:id", h.RemoveFromCart)
    r.GET("/deposit-request/:id", h.GetDepositRequest)

    if err := r.Run(":8080"); err != nil {
        log.Fatalf("failed to run server: %v", err)
    }
}