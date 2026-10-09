package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// DaySteps содержит данные о прогулке.
type DaySteps struct {
	Steps    int           // количество шагов.
	Duration time.Duration // длительность прогулки.
	personaldata.Personal
}

// Parse разбирает строку формата "678,0h50m" и заполняет поля DaySteps.
func (ds *DaySteps) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return fmt.Errorf("неверный формат данных %q: ожидается \"шаги,продолжительность\"", datastring)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}
	if steps <= 0 {
		return fmt.Errorf("некорректное количество шагов: %d", steps)
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}
	if duration <= 0 {
		return fmt.Errorf("некорректная продолжительность: %v", duration)
	}

	ds.Steps = steps
	ds.Duration = duration

	return nil
}

// ActionInfo возвращает информацию о прогулке.
func (ds DaySteps) ActionInfo() (string, error) {
	distance := spentenergy.Distance(ds.Steps, ds.Height)

	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories), nil
}
