package common

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

// datePatterns lists supported date formats in descending specificity.
var datePatterns = []string{
	`(\d{4}[/\-]\d{2}[/\-]\d{2})`,
}

var (
	dateRe   = regexp.MustCompile(`^` + datePatterns[0] + `\s+`)
	amountRe = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s+`)
	wordRe   = regexp.MustCompile(`^(\S+)\s*`)
)

// ParseTransactionText parses a natural-language transaction string and returns a
// partially- or fully-populated TransactionNoId plus human-readable warnings for
// fields that could not be resolved.
//
// Grammar (all tokens are whitespace-separated):
//
//	[DATE] AMOUNT CURRENCY [from ACCOUNT] [to ACCOUNT] [DESCRIPTION...]
//
// DATE     — YYYY/MM/DD or YYYY-MM-DD; defaults to today when omitted.
// AMOUNT   — decimal number (e.g. 100 or 50.5).
// CURRENCY — currency name, case-insensitive; fuzzy-matched.
// from/to  — keywords introducing account names; account name is matched
//
//	greedily against known accounts so that trailing description words
//	are not consumed into the account name.
func ParseTransactionText(
	text string,
	accounts []goserver.Account,
	currencies []goserver.Currency,
	today time.Time,
) (goserver.TransactionNoId, []string) {
	var warnings []string
	result := goserver.TransactionNoId{
		Date: today,
	}

	s := strings.TrimSpace(text)

	// --- optional date ---
	if m := dateRe.FindStringSubmatch(s); m != nil {
		raw := strings.ReplaceAll(m[1], "/", "-")
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			result.Date = t
		}
		s = s[len(m[0]):]
	}

	// --- amount ---
	amountStr := ""
	if m := amountRe.FindStringSubmatch(s); m != nil {
		amountStr = m[1]
		s = s[len(m[0]):]
	} else {
		warnings = append(warnings, "could not parse amount — expected a number")
		return result, warnings
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("invalid amount %q", amountStr))
		return result, warnings
	}

	// --- currency ---
	currencyID := ""
	if m := wordRe.FindStringSubmatch(s); m != nil {
		tok := m[1]
		id, partialWarn := fuzzyMatchCurrency(tok, currencies)
		if id != "" {
			currencyID = id
			if partialWarn != "" {
				warnings = append(warnings, partialWarn)
			}
		} else {
			warnings = append(warnings, fmt.Sprintf("currency %q not found", tok))
		}
		s = s[len(m[0]):]
	}

	// --- parse "from ACCOUNT" and "to ACCOUNT" clauses ---
	fromName, toName, remainder := parseFromToWithAccounts(s, accounts)

	var movements []goserver.Movement

	if fromName != "" {
		id, w := fuzzyMatchAccount(fromName, accounts)
		warnings = append(warnings, w...)
		movements = append(movements, goserver.Movement{
			AccountId:  id,
			CurrencyId: currencyID,
			Amount:     amount.Neg(),
		})
	}

	if toName != "" {
		id, w := fuzzyMatchAccount(toName, accounts)
		warnings = append(warnings, w...)
		movements = append(movements, goserver.Movement{
			AccountId:  id,
			CurrencyId: currencyID,
			Amount:     amount,
		})
	}

	// if neither from nor to, produce a single positive movement with no account
	if fromName == "" && toName == "" {
		warnings = append(warnings, "no 'from' or 'to' account found")
		movements = append(movements, goserver.Movement{
			CurrencyId: currencyID,
			Amount:     amount,
		})
	}

	result.Movements = movements

	desc := strings.TrimSpace(remainder)
	if desc != "" {
		result.Description = desc
	}

	return result, warnings
}

// parseFromTo splits the remaining string into optional fromName, toName,
// and a leftover description.
//
// For the final account in the string (the "to" clause when no further
// keywords follow) we use extractAccountName to avoid greedily consuming
// description words that happen to follow the account name.
func parseFromToWithAccounts(s string, accounts []goserver.Account) (fromName, toName, remainder string) {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)

	fromIdx := -1
	toIdx := -1

	if strings.HasPrefix(lower, "from ") {
		fromIdx = 0
	}
	if strings.HasPrefix(lower, "to ") {
		toIdx = 0
	}

	if fromIdx == 0 {
		rest := s[5:]
		lowerRest := strings.ToLower(rest)
		if idx := strings.Index(lowerRest, " to "); idx != -1 {
			toIdx = 5 + idx
		}
	}

	switch {
	case fromIdx == 0 && toIdx > 0:
		fromPart := s[5:toIdx]
		afterTo := s[toIdx+4:]
		toName, remainder = extractAccountName(afterTo, accounts)
		fromName = strings.TrimSpace(fromPart)

	case fromIdx == 0:
		rest := s[5:]
		fromName, remainder = extractAccountName(rest, accounts)

	case toIdx == 0:
		rest := s[3:]
		toName, remainder = extractAccountName(rest, accounts)

	default:
		remainder = s
	}

	return fromName, toName, remainder
}

// parseFromTo is the accounts-unaware variant kept for internal use.
func parseFromTo(s string) (fromName, toName, remainder string) {
	s = strings.TrimSpace(s)

	lower := strings.ToLower(s)

	fromIdx := -1
	toIdx := -1

	if strings.HasPrefix(lower, "from ") {
		fromIdx = 0
	}
	if strings.HasPrefix(lower, "to ") {
		toIdx = 0
	}

	// find "to " after a "from …" clause
	if fromIdx == 0 {
		rest := s[5:]
		lowerRest := strings.ToLower(rest)
		if idx := strings.Index(lowerRest, " to "); idx != -1 {
			toIdx = 5 + idx
		}
	}

	switch {
	case fromIdx == 0 && toIdx > 0:
		// "from X to Y …"
		fromPart := s[5:toIdx]
		afterTo := s[toIdx+4:] // skip " to "
		toName, remainder = splitFirstWordGroup(afterTo)
		fromName = strings.TrimSpace(fromPart)

	case fromIdx == 0:
		// "from X …"
		rest := s[5:]
		fromName, remainder = splitFirstWordGroup(rest)

	case toIdx == 0:
		// "to Y …"
		rest := s[3:]
		toName, remainder = splitFirstWordGroup(rest)

	default:
		remainder = s
	}

	return fromName, toName, remainder
}

// splitFirstWordGroup returns the first contiguous non-keyword words as the
// name, stopping when it encounters a standalone "to" or "from" keyword after
// at least one word has been consumed.
//
// Any words after the first unrecognized keyword become the remainder.
func splitFirstWordGroup(s string) (name, remainder string) {
	s = strings.TrimSpace(s)
	words := strings.Fields(s)
	var nameWords []string
	for i, w := range words {
		lw := strings.ToLower(w)
		if (lw == "to" || lw == "from") && i > 0 {
			remainder = strings.Join(words[i:], " ")
			break
		}
		nameWords = append(nameWords, w)
	}
	name = strings.Join(nameWords, " ")
	return name, remainder
}

// extractAccountName extracts the best matching account name from the beginning
// of s given a known list of accounts. It tries progressively shorter prefixes
// (greedy) so that "Others groceries" resolves to account "Others" with
// remainder "groceries" rather than failing to match "Others groceries".
//
// If no prefix matches any account, it returns the first keyword-terminated
// group as-is (fallback so we still produce a warning rather than silently
// dropping the text).
func extractAccountName(s string, accounts []goserver.Account) (name, remainder string) {
	words := strings.Fields(strings.TrimSpace(s))
	if len(words) == 0 {
		return "", ""
	}

	// Stop greedily expanding at "to"/"from" — these are reserved keywords.
	maxWords := 0
	for _, w := range words {
		lw := strings.ToLower(w)
		if (lw == "to" || lw == "from") && maxWords > 0 {
			break
		}
		maxWords++
	}

	// Try from longest to shortest prefix.
	for end := maxWords; end >= 1; end-- {
		candidate := strings.Join(words[:end], " ")
		low := strings.ToLower(candidate)
		for _, a := range accounts {
			if strings.ToLower(a.Name) == low || strings.Contains(strings.ToLower(a.Name), low) {
				return candidate, strings.TrimSpace(strings.Join(words[end:], " "))
			}
		}
	}

	// No match — return the keyword-terminated group so we can warn.
	name, remainder = splitFirstWordGroup(s)
	return name, remainder
}

// fuzzyMatchCurrency returns the currency ID for the best match (case-insensitive
// exact first, then partial). Empty string means not found.
func fuzzyMatchCurrency(token string, currencies []goserver.Currency) (id string, warn string) {
	low := strings.ToLower(token)
	// exact
	for _, c := range currencies {
		if strings.ToLower(c.Name) == low {
			return c.Id, ""
		}
	}
	// partial
	for _, c := range currencies {
		if strings.Contains(strings.ToLower(c.Name), low) {
			return c.Id, fmt.Sprintf("currency %q matched %q by partial name — please verify", token, c.Name)
		}
	}
	return "", ""
}

// fuzzyMatchAccount returns the account ID and any warnings.
func fuzzyMatchAccount(name string, accounts []goserver.Account) (id string, warnings []string) {
	low := strings.ToLower(name)
	// exact
	for _, a := range accounts {
		if strings.ToLower(a.Name) == low {
			return a.Id, nil
		}
	}
	// partial
	for _, a := range accounts {
		if strings.Contains(strings.ToLower(a.Name), low) {
			return a.Id, []string{fmt.Sprintf("%q matched %q by partial name — please verify", name, a.Name)}
		}
	}
	return "", []string{fmt.Sprintf("account %q not found", name)}
}
