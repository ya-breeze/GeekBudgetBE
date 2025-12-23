package common

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestMatchWithDetails(t *testing.T) {
	tests := []struct {
		name        string
		matcher     database.MatcherRuntime
		transaction goserver.Transaction
		wantResult  MatchResult
		wantMatched bool
	}{
		{
			name: "Match description",
			matcher: database.MatcherRuntime{
				DescriptionRegexp: regexp.MustCompile("^Uber.*"),
			},
			transaction: goserver.Transaction{
				Description: "Uber Eats",
			},
			wantResult:  MatchResultSuccess,
			wantMatched: true,
		},
		{
			name: "Miss match description",
			matcher: database.MatcherRuntime{
				DescriptionRegexp: regexp.MustCompile("^Uber.*"),
			},
			transaction: goserver.Transaction{
				Description: "Bolt",
			},
			wantResult:  MatchResultWrongDescription,
			wantMatched: false,
		},
		{
			name: "Match place",
			matcher: database.MatcherRuntime{
				PlaceRegexp: regexp.MustCompile(".*Prague.*"),
			},
			transaction: goserver.Transaction{
				Place: "Shell Prague",
			},
			wantResult:  MatchResultSuccess,
			wantMatched: true,
		},
		{
			name: "Miss match place",
			matcher: database.MatcherRuntime{
				PlaceRegexp: regexp.MustCompile(".*Prague.*"),
			},
			transaction: goserver.Transaction{
				Place: "Shell Berlin",
			},
			wantResult:  MatchResultWrongPartnerAccount, // Reusing existing error for now as per implementation
			wantMatched: false,
		},
		{
			name: "AliExpress False Positive Reproduction",
			matcher: database.MatcherRuntime{
				DescriptionRegexp: regexp.MustCompile("(?i)\\bAliExpress\\b"),
			},
			transaction: goserver.Transaction{
				Description: "Nákup: ",
				Place:       "MPLA S.R.O., PRAHA, CZ",
			},
			wantResult:  MatchResultWrongDescription,
			wantMatched: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchWithDetails(&tt.matcher, &tt.transaction)
			assert.Equal(t, tt.wantResult, got.Result)
			assert.Equal(t, tt.wantMatched, got.Matched)
		})
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name        string
		matcher     database.MatcherRuntime
		transaction goserver.Transaction
		wantResult  MatchResult
	}{
		{
			name: "AliExpress False Positive Reproduction (Match)",
			matcher: database.MatcherRuntime{
				DescriptionRegexp: regexp.MustCompile(".*"),
				PlaceRegexp:       regexp.MustCompile("(?i)\\bAliExpress\\b"),
			},
			transaction: goserver.Transaction{
				Description: "Nákup: ",
				Place:       "MPLA S.R.O., PRAHA, CZ",
			},
			wantResult: MatchResultWrongPartnerAccount, // Or whatever error corresponds to mismatch
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Match(&tt.matcher, &tt.transaction)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}
