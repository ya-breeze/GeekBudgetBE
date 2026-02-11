package common

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

// GetIncreases is a wrapper around utils.GetIncreases.
func GetIncreases(movements []goserver.Movement) map[string]decimal.Decimal {
	return utils.GetIncreases(movements)
}

// IsDuplicate is a wrapper around utils.IsDuplicate.
func IsDuplicate(t1Date time.Time, t1Movements []goserver.Movement, t2Date time.Time, t2Movements []goserver.Movement) bool {
	return utils.IsDuplicate(t1Date, t1Movements, t2Date, t2Movements)
}
