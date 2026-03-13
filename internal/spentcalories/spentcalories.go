package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

var (
	ErrInvalidWeight    = errors.New("неправильный вес")
	ErrInvalidHeight    = errors.New("неправильный рост")
	ErrInvalidFormat    = errors.New("неправильный формат")
	ErrNoSteps          = errors.New("нулевое количество шагов")
	ErrNegativeSteps    = errors.New("отрицательное число шагов")
	ErrNegativeDuration = errors.New("отрицательная продолжительность")
	ErrZeroDuration     = errors.New("нулевая продолжительность")
	ErrUnknownType      = errors.New("неизвестный тип тренировки")
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

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, ErrInvalidFormat
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	if err = validateStep(steps); err != nil {
		return 0, "", 0, err
	}

	action := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	if err = validateDuration(duration); err != nil {
		return 0, "", 0, err
	}

	return steps, action, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	speedKmh := dist / duration.Hours()
	return speedKmh
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, action, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	switch action {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", ErrUnknownType
	}
	speedKmh := meanSpeed(steps, height, duration)
	dist := distance(steps, height)
	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		action,
		duration.Hours(),
		dist,
		speedKmh,
		calories,
	), nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateStep(steps); err != nil {
		return 0, err
	}
	if err := validateDuration(duration); err != nil {
		return 0, err
	}
	if weight <= 0 {
		return 0, ErrInvalidWeight
	}
	if height <= 0 {
		return 0, ErrInvalidHeight
	}

	speedKmh := meanSpeed(steps, height, duration)
	calories := weight * speedKmh * duration.Minutes() / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateStep(steps); err != nil {
		return 0, err
	}
	if err := validateDuration(duration); err != nil {
		return 0, err
	}
	if weight <= 0 {
		return 0, ErrInvalidWeight
	}
	if height <= 0 {
		return 0, ErrInvalidHeight
	}

	speedKmh := meanSpeed(steps, height, duration)
	calories := (walkingCaloriesCoefficient * weight * speedKmh * duration.Minutes()) / minInH
	return calories, nil
}
