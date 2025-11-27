package handler

import (
    "net/http"
    "strconv"
    "strings"
    "time"

    "deposit-calculator/internal/app/repository"

    "github.com/gin-gonic/gin"
)

// Глобальная заявка/корзина (в памяти)
var depositRequest = &DepositRequest{
    ID:           1,
    Status:       "draft",
    MonthIDs:     []int{},           // выбранные ID месяцев
    Amounts:      make(map[int]int), // сумма по monthID
    Profits:      make(map[int]int), // прибыль по monthID
    TotalAmount:  0,
    TotalProfit:  0,
    CommonAmount: 100000,            // значение по умолчанию
}

type DepositRequest struct {
    ID           int
    Status       string
    MonthIDs     []int
    Amounts      map[int]int
    Profits      map[int]int
    TotalAmount  int
    TotalProfit  int
    CommonAmount int
}

type Handler struct {
    repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
    return &Handler{repo: repo}
}

// Предзаполнение корзины статическими данными (ID 1 и 2)
func (h *Handler) PrefillCart() {
    if len(depositRequest.MonthIDs) == 0 {
        ids := []int{1, 2}
        for _, mid := range ids {
            if _, err := h.repo.GetMonth(mid); err == nil {
                depositRequest.MonthIDs = append(depositRequest.MonthIDs, mid)
                depositRequest.Amounts[mid] = depositRequest.CommonAmount
                depositRequest.Profits[mid] = calculateProfit(h, mid)
            }
        }
        updateTotals()
    }
}

// GET / — главная (первые 3 карточки), поиск через ?query=
func (h *Handler) GetMonths(c *gin.Context) {
    query := c.Query("query")

    months, err := h.repo.GetMonths()
    if err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "cannot load months"})
        return
    }

    if query != "" {
        months, _ = h.repo.GetMonthsByName(query)
    }

    display := months
    if len(display) > 3 {
        display = display[:3]
    }

    c.HTML(http.StatusOK, "months_list.html", gin.H{
        "time":           time.Now().Format("15:04:05"),
        "months":         display,
        "depositRequest": depositRequest,
        "cartCount":      len(depositRequest.MonthIDs),
        "query":          query,
    })
}

// GET /month/:id — страница деталей месяца
func (h *Handler) GetMonthDetail(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "invalid month id"})
        return
    }

    m, err := h.repo.GetMonth(id)
    if err != nil {
        c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "month not found"})
        return
    }

    // Подстановки для страницы деталей
    m.Amount = depositRequest.CommonAmount
    m.Profit = calculateProfit(h, id)

    c.HTML(http.StatusOK, "month_detail.html", gin.H{
        "month":          m,
        "depositRequest": depositRequest,
        "cartCount":      len(depositRequest.MonthIDs),
    })
}

// POST /set-sum — задать общую сумму вклада
func (h *Handler) SetCommonAmount(c *gin.Context) {
    amountStr := c.PostForm("sum")
    if amountStr != "" {
        if v, err := strconv.Atoi(amountStr); err == nil && v >= 0 {
            depositRequest.CommonAmount = v
            // пересчёт для уже добавленных месяцев
            for _, mid := range depositRequest.MonthIDs {
                depositRequest.Amounts[mid] = depositRequest.CommonAmount
                depositRequest.Profits[mid] = calculateProfit(h, mid)
            }
            updateTotals()
        }
    }
    c.Redirect(http.StatusSeeOther, "/")
}

// POST /add-to-cart — form field monthId
func (h *Handler) AddToCart(c *gin.Context) {
    monthIDStr := c.PostForm("monthId")
    monthID, err := strconv.Atoi(monthIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month id"})
        return
    }

    // не добавлять повторно
    for _, id := range depositRequest.MonthIDs {
        if id == monthID {
            c.Redirect(http.StatusSeeOther, "/")
            return
        }
    }

    if _, err = h.repo.GetMonth(monthID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "month not found"})
        return
    }

    depositRequest.MonthIDs = append(depositRequest.MonthIDs, monthID)
    depositRequest.Amounts[monthID] = depositRequest.CommonAmount
    depositRequest.Profits[monthID] = calculateProfit(h, monthID)
    updateTotals()

    c.Redirect(http.StatusSeeOther, "/")
}

// POST /remove-from-cart/:id
func (h *Handler) RemoveFromCart(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    newIDs := []int{}
    for _, mid := range depositRequest.MonthIDs {
        if mid != id {
            newIDs = append(newIDs, mid)
        }
    }
    depositRequest.MonthIDs = newIDs
    delete(depositRequest.Amounts, id)
    delete(depositRequest.Profits, id)
    updateTotals()

    c.Redirect(http.StatusSeeOther, "/deposit-request/1")
}

// GET /deposit-request/:id
func (h *Handler) GetDepositRequest(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id != 1 {
        c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "request not found"})
        return
    }

    selected := []repository.Month{}
    for _, mid := range depositRequest.MonthIDs {
        m, err := h.repo.GetMonth(mid)
        if err == nil {
            m.Amount = depositRequest.CommonAmount
            m.Profit = depositRequest.Profits[mid]
            selected = append(selected, *m)
        }
    }

    c.HTML(http.StatusOK, "deposit_request.html", gin.H{
        "months":         selected,
        "depositRequest": depositRequest,
        "cartCount":      len(depositRequest.MonthIDs),
    })
}

// ====== вспомогательные функции ======

func calculateProfit(h *Handler, monthID int) int {
    month, err := h.repo.GetMonth(monthID)
    if err != nil {
        return 0
    }

    // "7.5% годовых" -> 7.5
    rateStr := strings.ReplaceAll(month.InterestRate, "годовых", "")
    rateStr = strings.TrimSpace(rateStr)
    rateStr = strings.TrimRight(rateStr, "%")
    rate, _ := strconv.ParseFloat(rateStr, 64)

    daysInYear := 365.0
    profit := float64(depositRequest.CommonAmount) * (rate / 100.0) * (float64(month.Days) / daysInYear)
    return int(profit)
}

func updateTotals() {
    depositRequest.TotalAmount = depositRequest.CommonAmount * len(depositRequest.MonthIDs)
    depositRequest.TotalProfit = 0
    for _, p := range depositRequest.Profits {
        depositRequest.TotalProfit += p
    }
}