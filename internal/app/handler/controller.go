package handler

import (
	"deposit-calculator/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetMonth(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid month ID"})
		return
	}

	month, err := h.Repository.GetMonth(id)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Month not found"})
		return
	}

	ctx.HTML(http.StatusOK, "month_detail.html", gin.H{
		"month": month,
	})
}

func (h *Handler) GetMonths(ctx *gin.Context) {
	var months []repository.Month
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		months, err = h.Repository.GetMonths()
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Internal server error"})
			return
		}
	} else {
		months, err = h.Repository.GetMonthsByName(searchQuery)
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Internal server error"})
			return
		}
	}

	calculationMonths := h.getCalculationMonths()
	calculationCount := len(calculationMonths)

	ctx.HTML(http.StatusOK, "months_list.html", gin.H{
		"time":             time.Now().Format("15:04:05"),
		"months":           months,
		"query":            searchQuery,
		"calculationCount": calculationCount,
	})
}

func (h *Handler) GetCalculationPage(ctx *gin.Context) {
	calculationMonths := h.getCalculationMonths()
	calculationCount := len(calculationMonths)

	ctx.HTML(http.StatusOK, "calculation.html", gin.H{
		"months":      calculationMonths,
		"monthsCount": calculationCount,
		"time":        time.Now().Format("15:04:05"),
	})
}

func (h *Handler) getCalculationMonths() []repository.Month {
	allMonths, err := h.Repository.GetMonths()
	if err != nil {
		return []repository.Month{}
	}
	return allMonths[:4]
}

func (h *Handler) GetCalculationCount(ctx *gin.Context) {
	calculationMonths := h.getCalculationMonths()
	calculationCount := len(calculationMonths)

	ctx.JSON(http.StatusOK, gin.H{
		"count": calculationCount,
	})
}
