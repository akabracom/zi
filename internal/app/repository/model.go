package repository

import (
	"fmt"
	"strings"
)

// Repository хранит одновременно slice и map
type Repository struct {
	months    []Month
	monthsMap map[int]Month
}

type Month struct {
	ID           int
	Name         string
	Days         int    // количество дней вклада (статично в этой лабе)
	InterestRate string // например "7.5% годовых"
	StartDate    string // статическая строка для отображения, например "02.02"
	EndDate      string // например "14.02"
	Season       string
	Image        string
	Category     string
	Description  string
	Holidays     string
	// Доп. поля для передачи в шаблон
	Amount int
	Profit int
}

func NewRepository() (*Repository, error) {
	r := &Repository{}
	r.monthsMap = make(map[int]Month)

	// Жёстко заданные данные (по требованию преподавателя — менять нельзя)
	m := []Month{
		{
			ID:           1,
			Name:         "Январь",
			Days:         31,
			InterestRate: "8% годовых",
			StartDate:    "01.01",
			EndDate:      "30.01",
			Season:       "Зима",
			Image:        "january.png",
			Category:     "Высокий сезон",
			Description:  "Начало года, период зимних праздников",
			Holidays:     "Новый год",
		},
		{
			ID:           2,
			Name:         "Февраль",
			Days:         28,
			InterestRate: "7.5% годовых",
			StartDate:    "02.02",
			EndDate:      "14.02",
			Season:       "Зима",
			Image:        "february.png",
			Category:     "Средний сезон",
			Description:  "Короткий месяц с повышенными ставками",
			Holidays:     "День защитника Отечества",
		},
		{
			ID:           3,
			Name:         "Март",
			Days:         31,
			InterestRate: "7.8% годовых",
			StartDate:    "10.03",
			EndDate:      "24.03",
			Season:       "Весна",
			Image:        "march.png",
			Category:     "Средний сезон",
			Description:  "Весеннее оживление банковских продуктов",
			Holidays:     "Международный женский день",
		},
		// Дополнительные месяцы остаются, но на главной показываем только первые 3
		{
			ID:           4,
			Name:         "Апрель",
			Days:         30,
			InterestRate: "6.5% годовых",
			StartDate:    "01.04",
			EndDate:      "30.04",
			Season:       "Весна",
			Image:        "april.png",
			Category:     "Низкий сезон",
			Description:  "Стабильный период",
			Holidays:     "День смеха",
		},
		{
			ID:           5,
			Name:         "Май",
			Days:         31,
			InterestRate: "6.7% годовых",
			StartDate:    "01.05",
			EndDate:      "31.05",
			Season:       "Весна",
			Image:        "may.png",
			Category:     "Средний сезон",
			Description:  "Месяц майских праздников",
			Holidays:     "День Победы",
		},
	}

	r.months = m
	for _, mm := range m {
		r.monthsMap[mm.ID] = mm
	}

	return r, nil
}

func (r *Repository) GetMonths() ([]Month, error) {
	if len(r.months) == 0 {
		return nil, fmt.Errorf("массив месяцев пустой")
	}
	return r.months, nil
}

func (r *Repository) GetMonthsByName(query string) ([]Month, error) {
	if query == "" {
		return r.GetMonths()
	}
	var filtered []Month
	for _, m := range r.months {
		if containsIgnoreCase(m.Name, query) {
			filtered = append(filtered, m)
		}
	}
	return filtered, nil
}

func (r *Repository) GetMonth(id int) (*Month, error) {
	if m, ok := r.monthsMap[id]; ok {
		tmp := m
		return &tmp, nil
	}
	return nil, fmt.Errorf("month with ID %d not found", id)
}

func (r *Repository) GetMonthsMap() map[int]Month {
	// возвращаем копию чтобы никто не менял исходный map
	out := make(map[int]Month, len(r.monthsMap))
	for k, v := range r.monthsMap {
		out[k] = v
	}
	return out
}

func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}
