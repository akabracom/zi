package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Month struct {
	ID           int
	Name         string // Название месяца
	Description  string // Описание особенностей месяца
	Days         int    // Количество дней
	InterestRate string // Процентная ставка
	Season       string // Сезонность
	Holidays     string // Праздничные дни
	Image        string // Путь к изображению
	Category     string // Категория
}

func (r *Repository) GetMonths() ([]Month, error) {
	months := []Month{
		{
			ID:           1,
			Name:         "Январь",
			Description:  "Начало года, период зимних праздников и высоких ставок по вкладам",
			Days:         31,
			InterestRate: "7.5% годовых",
			Season:       "Зима",
			Holidays:     "Новый год, Рождество",
			Image:        "january.png",
			Category:     "Высокий сезон",
		},
		{
			ID:           2,
			Name:         "Февраль",
			Description:  "Короткий месяц с повышенными ставками для привлечения клиентов",
			Days:         28,
			InterestRate: "7.2% годовых",
			Season:       "Зима",
			Holidays:     "День защитника Отечества",
			Image:        "february.png",
			Category:     "Средний сезон",
		},
		{
			ID:           3,
			Name:         "Март",
			Description:  "Весеннее оживление банковских продуктов и стабилизация ставок",
			Days:         31,
			InterestRate: "6.8% годовых",
			Season:       "Весна",
			Holidays:     "Международный женский день",
			Image:        "march.png",
			Category:     "Средний сезон",
		},
		{
			ID:           4,
			Name:         "Апрель",
			Description:  "Стабильный период с умеренными процентными ставками",
			Days:         30,
			InterestRate: "6.5% годовых",
			Season:       "Весна",
			Holidays:     "День смеха",
			Image:        "april.png",
			Category:     "Низкий сезон",
		},
		{
			ID:           5,
			Name:         "Май",
			Description:  "Месяц майских праздников с особыми условиями по вкладам",
			Days:         31,
			InterestRate: "6.7% годовых",
			Season:       "Весна",
			Holidays:     "День Победы, Праздник весны и труда",
			Image:        "may.png",
			Category:     "Средний сезон",
		},
	}

	if len(months) == 0 {
		return nil, fmt.Errorf("массив месяцев пустой")
	}
	return months, nil
}

func (r *Repository) GetMonthsByName(query string) ([]Month, error) {
	allMonths, err := r.GetMonths()
	if err != nil {
		return nil, err
	}

	var filtered []Month
	for _, month := range allMonths {
		if containsIgnoreCase(month.Name, query) {
			filtered = append(filtered, month)
		}
	}

	return filtered, nil
}

func (r *Repository) GetMonth(id int) (*Month, error) {
	months, err := r.GetMonths()
	if err != nil {
		return nil, err
	}

	for _, month := range months {
		if month.ID == id {
			return &month, nil
		}
	}

	return nil, fmt.Errorf("month with ID %d not found", id)
}

// Вспомогательная функция для поиска без учета регистра
func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}
