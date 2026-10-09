package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Виды тренировок.
const (
	walking = "Ходьба"
	running = "Бег"
)

// ErrUnknownTraining возвращается при неизвестном типе тренировки.
var ErrUnknownTraining = errors.New("неизвестный тип тренировки")

// Training содержит данные о тренировке.
type Training struct {
	Steps        int           // количество шагов.
	TrainingType string        // тип тренировки (бег или ходьба).
	Duration     time.Duration // длительность тренировки.
	personaldata.Personal
}

// Parse разбирает строку формата "3456,Ходьба,3h00m" и заполняет поля Training.
func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return fmt.Errorf("неверный формат данных %q: ожидается \"шаги,вид активности,продолжительность\"", datastring)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}
	if steps <= 0 {
		return fmt.Errorf("некорректное количество шагов: %d", steps)
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}
	if duration <= 0 {
		return fmt.Errorf("некорректная продолжительность: %v", duration)
	}

	t.Steps = steps
	t.TrainingType = parts[1]
	t.Duration = duration

	return nil
}

// ActionInfo возвращает информацию о тренировке.
func (t Training) ActionInfo() (string, error) {
	distance := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var (
		calories float64
		err      error
	)
	switch t.TrainingType {
	case running:
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case walking:
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", ErrUnknownTraining
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, speed, calories), nil
}
