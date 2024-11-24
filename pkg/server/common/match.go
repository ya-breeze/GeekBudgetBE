package common

import (
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type MatchResult int

const (
	MatchResultSuccess MatchResult = iota
	MatchResultWrongDescription
	MatchResultWrongPartnerAccount
)

func Match(matcher *database.MatcherRuntime, transaction *goserver.Transaction) MatchResult {
	if matcher.DescriptionRegexp != nil && !matcher.DescriptionRegexp.MatchString(transaction.Description) {
		return MatchResultWrongDescription
	}

	if matcher.PartnerAccountRegexp != nil &&
		!matcher.PartnerAccountRegexp.MatchString(transaction.PartnerAccount) {
		return MatchResultWrongPartnerAccount
	}

	return MatchResultSuccess
}
