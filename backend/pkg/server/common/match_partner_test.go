package common

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestMatchPartnerName(t *testing.T) {
	tests := []struct {
		name        string
		matcher     database.MatcherRuntime
		transaction goserver.Transaction
		wantResult  MatchResult
		wantMatched bool
	}{
		{
			name: "Match partner name",
			matcher: database.MatcherRuntime{
				PartnerNameRegexp: regexp.MustCompile("^Lidl$"),
			},
			transaction: goserver.Transaction{
				PartnerName: "Lidl",
			},
			wantResult:  MatchResultSuccess,
			wantMatched: true,
		},
		{
			name: "Miss match partner name",
			matcher: database.MatcherRuntime{
				PartnerNameRegexp: regexp.MustCompile("^Lidl$"),
			},
			transaction: goserver.Transaction{
				PartnerName: "Tesco",
			},
			wantResult:  MatchResultWrongPartnerName, // This constant needs to be added
			wantMatched: false,
		},
		{
			name: "Match partner name with alternation (Multiple Vendors)",
			matcher: database.MatcherRuntime{
				PartnerNameRegexp: regexp.MustCompile("(?i)(Lidl|Albert|Billa)"),
			},
			transaction: goserver.Transaction{
				PartnerName: "Albert",
			},
			wantResult:  MatchResultSuccess,
			wantMatched: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchWithDetails(&tt.matcher, &tt.transaction)
			// Until we implement the fix, these are expected to fail if the logic ignores PartnerNameRegexp
			// However, in the current broken state, MatchWithDetails returns MatchResultSuccess if checked fields are nil/empty.
			// Since PartnerNameRegexp is not checked, it will likely return Success even for mismatches (False Positive)
			// OR validation logic might be totally missing.

			// For the purpose of "Reproduction", we EXPECT MatchResultSuccess if it's currently ignored.
			// But for TDD, we write what we WANT.
			// I will assert what we WANT, so the test fails, proving the bug/missing feature.

			// Note: MatchResultWrongPartnerName constant doesn't exist yet, so this code won't compile without it.
			// I should add the constant first or use a placeholder in the test and update it.
			// But wait, I can't add the constant in the test file efficiently without modifying source code.
			// I will use MatchResultSuccess for the "Miss match" case temporarily to see it PASS (meaning it falsely matches)
			// NO, that's confusing.
			// I'll assume I'll add the constant in the next step.
			// For initial "fail", I expect "MatchResultWrongPartnetName" or "Success" depending on perspective.

			// Let's first run it and see it fail compilation or assertion.
			// To make it compile, I'll comment out the specific Result check for now and rely on Matched bool.

			assert.Equal(t, tt.wantMatched, got.Matched)
		})
	}
}
