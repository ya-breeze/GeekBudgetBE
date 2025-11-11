package common

import (
	"fmt"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type MatchResult int

const (
	MatchResultSuccess MatchResult = iota
	MatchResultWrongDescription
	MatchResultWrongPartnerAccount
)

// MatchDetails contains detailed information about why a matcher matched or didn't match a transaction
type MatchDetails struct {
	Result                MatchResult
	Matched               bool
	FailureReason         string
	DescriptionMatched    bool
	PartnerAccountMatched bool
}

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

// MatchWithDetails returns detailed information about the match result
func MatchWithDetails(matcher *database.MatcherRuntime, transaction *goserver.Transaction) MatchDetails {
	details := MatchDetails{
		DescriptionMatched:    true,
		PartnerAccountMatched: true,
	}

	// Check description regex
	if matcher.DescriptionRegexp != nil {
		if !matcher.DescriptionRegexp.MatchString(transaction.Description) {
			details.Result = MatchResultWrongDescription
			details.FailureReason = fmt.Sprintf(
				"Description regex %q doesn't match transaction description %q",
				matcher.DescriptionRegexp.String(),
				transaction.Description,
			)
			details.DescriptionMatched = false
			details.Matched = false
			return details
		}
	}

	// Check partner account regex
	if matcher.PartnerAccountRegexp != nil {
		if !matcher.PartnerAccountRegexp.MatchString(transaction.PartnerAccount) {
			details.Result = MatchResultWrongPartnerAccount
			details.FailureReason = fmt.Sprintf(
				"Partner account regex %q doesn't match transaction partner account %q",
				matcher.PartnerAccountRegexp.String(),
				transaction.PartnerAccount,
			)
			details.PartnerAccountMatched = false
			details.Matched = false
			return details
		}
	}

	// All checks passed
	details.Result = MatchResultSuccess
	details.Matched = true
	details.FailureReason = ""
	return details
}
