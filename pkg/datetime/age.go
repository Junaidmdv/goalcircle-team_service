package datetime

import (
	"time"
)

type DateCalculator interface {
	CalculateAge(time.Time) int32
}

type dateCalculator struct{}

func NewDateCalculator() DateCalculator {
	return &dateCalculator{}
}

func (dc *dateCalculator) CalculateAge(dob time.Time) int32 {
	now := time.Now()

	if dob.After(now) {
		return 0
	}

	age := now.Year() - dob.Year()

	if now.Month() < dob.Month() ||
		(now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}

	return int32(age)
}
