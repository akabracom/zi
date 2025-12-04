package repository

import (
    "errors"
    "strconv"
    "strings"
    "time"

    "di_web/internal/models"
    "gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

// ====================== USER / DRAFT ======================

func (r *Repository) FindDraftWithItems(email string) (*models.DepositApplication, error) {
    u, err := r.ensureUser(email)
    if err != nil {
        return nil, err
    }

    var app models.DepositApplication
    err = r.db.Where("user_id = ? AND status = 'draft'", u.ID).
        Preload("Items").Preload("Items.Offer").First(&app).Error
    if err == gorm.ErrRecordNotFound {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &app, nil
}

func (r *Repository) ensureUser(email string) (*models.User, error) {
    var u models.User
    err := r.db.Where("email = ?", email).First(&u).Error
    if err == gorm.ErrRecordNotFound {
        u = models.User{Email: email, Name: "Student", CreatedAt: time.Now()}
        if err = r.db.Create(&u).Error; err != nil {
            return nil, err
        }
        return &u, nil
    }
    return &u, err
}

func (r *Repository) GetOrCreateDraft(email string) (*models.DepositApplication, error) {
    u, err := r.ensureUser(email)
    if err != nil {
        return nil, err
    }

    var app models.DepositApplication
    err = r.db.Where("user_id = ? AND status = 'draft'", u.ID).First(&app).Error
    if err == gorm.ErrRecordNotFound {
        app = models.DepositApplication{
            UserID:       u.ID,
            Status:       "draft",
            CreatedAt:    time.Now(),
            UpdatedAt:    time.Now(),
            CommonAmount: 100000,
            TotalAmount:  0,
            TotalProfit:  0,
        }
        if err = r.db.Create(&app).Error; err != nil {
            return nil, err
        }
        return &app, nil
    }
    return &app, err
}

// ====================== OFFERS ======================

func (r *Repository) GetDepositOffers(query string) ([]models.DepositOffer, error) {
    var out []models.DepositOffer
    q := r.db.Model(&models.DepositOffer{}).Where("is_active = true")
    if query != "" {
        like := "%" + strings.ToLower(query) + "%"
        q = q.Where("lower(name) LIKE ? OR lower(description) LIKE ?", like, like)
    }
    err := q.Order("name ASC").Find(&out).Error
    return out, err
}

func (r *Repository) GetOfferByID(id int) (*models.DepositOffer, error) {
    var o models.DepositOffer
    if err := r.db.First(&o, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.New("not found")
        }
        return nil, err
    }
    return &o, nil
}

// ====================== CALC ======================

func (r *Repository) CalculateProfit(ratePct float64, days int, amount float64) float64 {
    return amount * (ratePct/100.0) * (float64(days)/365.0)
}

func (r *Repository) RecalculateTotals(appID int64) error {
    var app models.DepositApplication
    if err := r.db.First(&app, appID).Error; err != nil {
        return err
    }

    var items []models.DepositApplicationItem
    if err := r.db.Where("application_id = ?", app.ID).Find(&items).Error; err != nil {
        return err
    }

    totalAmount := app.CommonAmount * float64(len(items))
    totalProfit := 0.0
    for _, it := range items {
        totalProfit += it.Profit
    }

    return r.db.Model(&app).Updates(map[string]any{
        "total_amount": totalAmount,
        "total_profit": totalProfit,
        "updated_at":   time.Now(),
    }).Error
}

// ====================== DRAFT ITEMS ======================

func (r *Repository) AddOfferToDraft(email string, offerID int) error {
    app, err := r.GetOrCreateDraft(email)
    if err != nil {
        return err
    }

    var off models.DepositOffer
    if err := r.db.First(&off, offerID).Error; err != nil {
        return err
    }

    var item models.DepositApplicationItem
    err = r.db.Where("application_id = ? AND offer_id = ?", app.ID, off.ID).
        First(&item).Error
    if err == gorm.ErrRecordNotFound {
        profit := r.CalculateProfit(off.InterestRatePct, off.Days, app.CommonAmount)
        item = models.DepositApplicationItem{
            ApplicationID: app.ID,
            OfferID:       off.ID,
            Quantity:      1,
            Amount:        app.CommonAmount,
            Profit:        profit,
        }
        if err := r.db.Create(&item).Error; err != nil {
            return err
        }
        return r.RecalculateTotals(app.ID)
    }
    return err
}

func (r *Repository) RemoveOfferFromDraft(email string, offerID int) error {
    app, err := r.GetOrCreateDraft(email)
    if err != nil {
        return err
    }

    res := r.db.Where("application_id = ? AND offer_id = ?", app.ID, offerID).
        Delete(&models.DepositApplicationItem{})
    if res.Error != nil {
        return res.Error
    }
    return r.RecalculateTotals(app.ID)
}

func (r *Repository) LoadDraftWithItems(email string) (*models.DepositApplication, error) {
    app, err := r.GetOrCreateDraft(email)
    if err != nil {
        return nil, err
    }
    if err := r.db.Where("id = ?", app.ID).
        Preload("Items").Preload("Items.Offer").First(app).Error; err != nil {
        return nil, err
    }
    return app, nil
}

// ====================== APPLICATIONS ======================

func (r *Repository) GetApplicationByID(id int64) (*models.DepositApplication, error) {
    var app models.DepositApplication
    err := r.db.Where("id = ? AND status != 'deleted'", id).
        Preload("Items").Preload("Items.Offer").First(&app).Error
    if err == gorm.ErrRecordNotFound {
        return nil, errors.New("not found")
    }
    return &app, err
}

func (r *Repository) DeleteApplicationLogical(id int64) error {
    return r.db.Model(&models.DepositApplication{}).
        Where("id = ?", id).
        Updates(map[string]any{
            "status":     "deleted",
            "updated_at": time.Now(),
        }).Error
}

// ====================== COMMON AMOUNT ======================

// обновление общей суммы по email (для главной)
func (r *Repository) UpdateCommonAmount(email string, amount float64) error {
    app, err := r.GetOrCreateDraft(email)
    if err != nil {
        return err
    }
    return r.UpdateCommonAmountByAppID(app.ID, amount)
}

// обновление общей суммы по ID заявки (для SaveRequest)
func (r *Repository) UpdateCommonAmountByAppID(appID int64, amount float64) error {
    var app models.DepositApplication
    if err := r.db.First(&app, appID).Error; err != nil {
        return err
    }

    if err := r.db.Model(&app).Updates(map[string]any{
        "common_amount": amount,
        "updated_at":    time.Now(),
    }).Error; err != nil {
        return err
    }

    var items []models.DepositApplicationItem
    if err := r.db.Where("application_id = ?", app.ID).
        Preload("Offer").Find(&items).Error; err != nil {
        return err
    }

    for _, it := range items {
        days := it.Offer.Days
        profit := r.CalculateProfit(it.Offer.InterestRatePct, days, amount)

        if err := r.db.Model(&it).Updates(map[string]any{
            "amount": amount,
            "profit": profit,
        }).Error; err != nil {
            return err
        }
    }

    return r.RecalculateTotals(app.ID)
}

// ====================== DATES ======================

func (r *Repository) UpdateItemDates(appID int64, offerID int, start, end string) error {
    var item models.DepositApplicationItem
    if err := r.db.Where("application_id = ? AND offer_id = ?", appID, offerID).
        Preload("Offer").First(&item).Error; err != nil {
        return err
    }

    item.StartDate = start
    item.EndDate = end

    if start != "" && end != "" {
        s, serr := strconv.Atoi(start)
        e, eerr := strconv.Atoi(end)
        if serr == nil && eerr == nil && e >= s {
            days := e - s + 1
            profit := r.CalculateProfit(item.Offer.InterestRatePct, days, item.Amount)
            item.Profit = profit
        }
    }

    if err := r.db.Save(&item).Error; err != nil {
        return err
    }

    return r.RecalculateTotals(appID)
}
