package handler

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "di_web/internal/app/repository"
    "di_web/internal/models"

    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
)

// ====================== HANDLER STRUCT ======================

type Handler struct {
    repo  *repository.Repository
    minio *minio.Client
}

func New(repo *repository.Repository, m *minio.Client) *Handler {
    return &Handler{repo: repo, minio: m}
}

// ====================== VIEW MODELS ======================

type MonthVM struct {
    ID           int
    Name         string
    Days         int
    InterestRate string
    StartDate    string // день начала (строка)
    EndDate      string // день конца (строка)
    Season       string
    Image        string
    Category     string
    Description  string
    Holidays     string
    Amount       int
    Profit       int
}

type DepositRequestVM struct {
    ID           int64
    CommonAmount int
    TotalAmount  int
    TotalProfit  int
}

func toVMRate(ratePct float64) string {
    return fmt.Sprintf("%.2f%% годовых", ratePct)
}

// ====================== HTML HANDLERS ======================

// список всех вкладов
func (h *Handler) GetMonths(c *gin.Context) {
    query := c.Query("query")

    offers, err := h.repo.GetDepositOffers(query)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Не удалось загрузить вклады"})
        return
    }

    draft, _ := h.repo.LoadDraftWithItems(currentEmail(c))

    vms := make([]MonthVM, 0, len(offers))
    for _, o := range offers {
        vms = append(vms, MonthVM{
            ID:           o.ID,
            Name:         o.Name,
            Days:         o.Days,
            InterestRate: toVMRate(o.InterestRatePct),
            StartDate:    o.StartDate,
            EndDate:      o.EndDate,
            Season:       o.Season,
            Image:        o.ImageKey,
            Category:     o.Category,
            Description:  o.Description,
            Holidays:     o.Holidays,
        })
    }

    c.HTML(http.StatusOK, "months_list.html", gin.H{
        "months":         vms,
        "depositRequest": toReqVM(draft),
        "cartCount":      len(draft.Items),
        "query":          query,
    })
}

// детали одного вклада
func (h *Handler) GetMonthDetail(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID"})
        return
    }

    o, err := h.repo.GetOfferByID(id)
    if err != nil {
        c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Вклад не найден"})
        return
    }

    draft, _ := h.repo.LoadDraftWithItems(currentEmail(c))
    profit := int(h.repo.CalculateProfit(o.InterestRatePct, o.Days, draft.CommonAmount))

    vm := MonthVM{
        ID:           o.ID,
        Name:         o.Name,
        Days:         o.Days,
        InterestRate: toVMRate(o.InterestRatePct),
        StartDate:    o.StartDate,
        EndDate:      o.EndDate,
        Season:       o.Season,
        Image:        o.ImageKey,
        Category:     o.Category,
        Description:  o.Description,
        Holidays:     o.Holidays,
        Amount:       int(draft.CommonAmount),
        Profit:       profit,
    }

    c.HTML(http.StatusOK, "month_detail.html", gin.H{
        "month":         vm,
        "depositRequest": toReqVM(draft),
        "cartCount":      len(draft.Items),
    })
}

// добавить вклад в черновик
func (h *Handler) AddToCart(c *gin.Context) {
    offerID, err := strconv.Atoi(c.PostForm("offer_id"))
    if err != nil || offerID <= 0 {
        c.Redirect(http.StatusSeeOther, "/")
        return
    }
    _ = h.repo.AddOfferToDraft(currentEmail(c), offerID)
    c.Redirect(http.StatusSeeOther, "/")
}

// удалить вклад из черновика
func (h *Handler) RemoveFromCart(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.Redirect(http.StatusSeeOther, "/")
        return
    }
    _ = h.repo.RemoveOfferFromDraft(currentEmail(c), id)
    draft, _ := h.repo.LoadDraftWithItems(currentEmail(c))
    c.Redirect(http.StatusSeeOther, fmt.Sprintf("/deposit-request/%d", draft.ID))
}

// страница заявки
func (h *Handler) GetDepositRequest(c *gin.Context) {
    id64, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Заявка не найдена"})
        return
    }

    app, err := h.repo.GetApplicationByID(id64)
    if err != nil {
        c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Заявка не найдена"})
        return
    }

    vms := make([]MonthVM, 0, len(app.Items))
    for _, it := range app.Items {
        vms = append(vms, MonthVM{
            ID:           it.OfferID,
            Name:         it.Offer.Name,
            Days:         it.Offer.Days,
            InterestRate: toVMRate(it.Offer.InterestRatePct),
            StartDate:    it.StartDate,
            EndDate:      it.EndDate,
            Season:       it.Offer.Season,
            Image:        it.Offer.ImageKey,
            Category:     it.Offer.Category,
            Description:  it.Offer.Description,
            Holidays:     it.Offer.Holidays,
            Amount:       int(it.Amount),
            Profit:       int(it.Profit),
        })
    }

    c.HTML(http.StatusOK, "deposit_request.html", gin.H{
        "months":         vms,
        "depositRequest": toReqVM(app),
        "cartCount":      len(app.Items),
    })
}

// сохранить сумму и даты по заявке (одна кнопка)
func (h *Handler) SaveRequest(c *gin.Context) {
    id64, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID заявки"})
        return
    }

    // 1) обновляем общую сумму по ID заявки
    amountStr := strings.TrimSpace(c.PostForm("sum"))
    if amountStr != "" {
        if v, err := strconv.Atoi(amountStr); err == nil && v >= 0 {
            _ = h.repo.UpdateCommonAmountByAppID(id64, float64(v))
        }
    }

    // 2) обновляем дни и прибыль только для первого и последнего элемента
    app, err := h.repo.GetApplicationByID(id64)
    if err != nil {
        c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Заявка не найдена"})
        return
    }

    if len(app.Items) > 0 {
        // Первый элемент
        firstItem := app.Items[0]
        startKey := fmt.Sprintf("start_%d", firstItem.OfferID)
        endKey := fmt.Sprintf("end_%d", firstItem.OfferID)
        start := strings.TrimSpace(c.PostForm(startKey))
        end := strings.TrimSpace(c.PostForm(endKey))
        _ = h.repo.UpdateItemDates(app.ID, firstItem.OfferID, start, end)

        // Последний элемент (если он отличается от первого)
        if len(app.Items) > 1 {
            lastItem := app.Items[len(app.Items)-1]
            startKey = fmt.Sprintf("start_%d", lastItem.OfferID)
            endKey = fmt.Sprintf("end_%d", lastItem.OfferID)
            start = strings.TrimSpace(c.PostForm(startKey))
            end = strings.TrimSpace(c.PostForm(endKey))
            _ = h.repo.UpdateItemDates(app.ID, lastItem.OfferID, start, end)
        }
    }

    c.Redirect(http.StatusSeeOther, fmt.Sprintf("/deposit-request/%d", id64))
}

// ====================== REST API ======================

func (h *Handler) ListOffers(c *gin.Context) {
    items, err := h.repo.GetDepositOffers(c.Query("query"))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
        return
    }
    c.JSON(http.StatusOK, items)
}

func (h *Handler) GetDraft(c *gin.Context) {
    app, err := h.repo.LoadDraftWithItems(currentEmail(c))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load"})
        return
    }
    c.JSON(http.StatusOK, app)
}

func (h *Handler) GetAppByID(c *gin.Context) {
    id64, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    app, err := h.repo.GetApplicationByID(id64)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    c.JSON(http.StatusOK, app)
}

func (h *Handler) AddOffer(c *gin.Context) {
    var body struct {
        OfferID int `json:"offer_id"`
    }
    if err := c.ShouldBindJSON(&body); err != nil || body.OfferID <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    if err := h.repo.AddOfferToDraft(currentEmail(c), body.OfferID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "cannot add"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) DeleteApp(c *gin.Context) {
    id64, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    if err := h.repo.DeleteApplicationLogical(id64); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// ====================== HELPERS ======================

func currentEmail(c *gin.Context) string {
    if v := c.GetHeader("X-User-Email"); v != "" {
        return v
    }
    return "student@example.com"
}

func toReqVM(app *models.DepositApplication) DepositRequestVM {
    if app == nil {
        return DepositRequestVM{ID: 0, CommonAmount: 100000}
    }
    return DepositRequestVM{
        ID:           app.ID,
        CommonAmount: int(app.CommonAmount),
        TotalAmount:  int(app.TotalAmount),
        TotalProfit:  int(app.TotalProfit),
    }
}
