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
	MatchResultWrongPartnerName
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

	if matcher.PlaceRegexp != nil &&
		!matcher.PlaceRegexp.MatchString(transaction.Place) {
		return MatchResultWrongPartnerAccount // Using existing error code or add a new one if strictly necessary, but standard MatchResult structure might suffice for boolean check
	}

	if matcher.PartnerNameRegexp != nil &&
		!matcher.PartnerNameRegexp.MatchString(transaction.PartnerName) {
		return MatchResultWrongPartnerName
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

	// Check partner name regex
	if matcher.PartnerNameRegexp != nil {
		if !matcher.PartnerNameRegexp.MatchString(transaction.PartnerName) {
			details.Result = MatchResultWrongPartnerName
			details.FailureReason = fmt.Sprintf(
				"Partner name regex %q doesn't match transaction partner name %q",
				matcher.PartnerNameRegexp.String(),
				transaction.PartnerName,
			)
			details.Matched = false
			return details
		}
	}

	// Check place regex
	if matcher.PlaceRegexp != nil {
		if !matcher.PlaceRegexp.MatchString(transaction.Place) {
			details.Result = MatchResultWrongPartnerAccount // Reusing existing error type or could define new one
			details.FailureReason = fmt.Sprintf(
				"Place regex %q doesn't match transaction place %q",
				matcher.PlaceRegexp.String(),
				transaction.Place,
			)
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
