package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

var (
	ErrInvalidFormat    = errors.New("неправильный формат")
	ErrNoSteps          = errors.New("нулевое количество шагов")
	ErrNegativeSteps    = errors.New("отрицательное число шагов")
	ErrNegativeDuration = errors.New("отрицательная продолжительность")
	ErrZeroDuration     = errors.New("нулевая продолжительность")
)

func validateStep(steps int) error {
	if steps == 0 {
		return ErrNoSteps
	}
	if steps < 0 {
		return ErrNegativeSteps
	}
	return nil
}

func validateDuration(duration time.Duration) error {
	if duration < 0 {
		return ErrNegativeDuration
	}
	if duration == 0 {
		return ErrZeroDuration
	}
	return nil
}

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, ErrInvalidFormat
	}
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	if err = validateStep(steps); err != nil {
		return 0, 0, err
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	if err = validateDuration(duration); err != nil {
		return 0, 0, err
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("error:", err)
		return ""
	}
	distanceKm := (float64(steps) * stepLength) / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		log.Println("error:", err)
		return ""
	}
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)
}
