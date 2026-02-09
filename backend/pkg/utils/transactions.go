package utils

import (
	"math"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

// GetIncreases calculates the net increase per currency from a list of movements.
func GetIncreases(movements []goserver.Movement) map[string]float64 {
	pos := make(map[string]float64)
	neg := make(map[string]float64)
	for _, m := range movements {
		if m.Amount > 0 {
			pos[m.CurrencyId] += m.Amount
		} else {
			neg[m.CurrencyId] -= m.Amount
		}
	}

	res := make(map[string]float64)
	for c, p := range pos {
		n := neg[c]
		if p > n {
			res[c] = p
		} else {
			res[c] = n
		}
	}
	for c, n := range neg {
		if _, ok := res[c]; !ok {
			res[c] = n
		}
	}
	return res
}

// IsDuplicate checks if two transactions are likely duplicates based on date and amounts.
func IsDuplicate(t1Date time.Time, t1Movements []goserver.Movement, t2Date time.Time, t2Movements []goserver.Movement) bool {
	// 1. Time check (+/- 2 days)
	delta := t1Date.Sub(t2Date)
	if delta < 0 {
		delta = -delta
	}
	if delta > 2*time.Hour*24 {
		return false
	}

	// 2. Amount check (sum of increases per currency)
	inc1 := GetIncreases(t1Movements)
	inc2 := GetIncreases(t2Movements)

	if len(inc1) != len(inc2) {
		return false
	}

	for c, v1 := range inc1 {
		v2, ok := inc2[c]
		if !ok || math.Abs(v1-v2) >= 0.01 {
			return false
		}
	}

	return true
}
