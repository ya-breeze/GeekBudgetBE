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
				DescriptionRegexp: regexp.MustCompile(`(?i)\bAliExpress\b`),
			},
			transaction: goserver.Transaction{
				Description: "Nákup: ",
				Place:       "MPLA S.R.O., PRAHA, CZ",
			},
			wantResult:  MatchResultWrongDescription,
			wantMatched: false,
		},
		{
			name: "Simplified match: case-insensitive",
			matcher: database.MatcherRuntime{
				Matcher:        &goserver.Matcher{Simplified: true},
				Keywords:       []string{"Uber"},
				KeywordOutputs: []string{"Uber"},
				KeywordRegexps: []*regexp.Regexp{
					regexp.MustCompile(`(?i)\bUber\b`),
				},
			},
			transaction: goserver.Transaction{
				Description: "uber eats",
			},
			wantResult:  MatchResultSuccess,
			wantMatched: true,
		},
		{
			name: "Simplified match: whole word only",
			matcher: database.MatcherRuntime{
				Matcher:        &goserver.Matcher{Simplified: true},
				Keywords:       []string{"Uber"},
				KeywordOutputs: []string{"Uber"},
				KeywordRegexps: []*regexp.Regexp{
					regexp.MustCompile(`(?i)\bUber\b`),
				},
			},
			transaction: goserver.Transaction{
				Description: "Uberrimo",
			},
			wantResult:  MatchResultWrongDescription,
			wantMatched: false,
		},
		{
			name: "Simplified match: first keyword wins",
			matcher: database.MatcherRuntime{
				Matcher:        &goserver.Matcher{Simplified: true},
				Keywords:       []string{"Uber", "Eats"},
				KeywordOutputs: []string{"Uber", "Eats"},
				KeywordRegexps: []*regexp.Regexp{
					regexp.MustCompile(`(?i)\bUber\b`),
					regexp.MustCompile(`(?i)\bEats\b`),
				},
			},
			transaction: goserver.Transaction{
				Description: "Uber Eats",
			},
			wantResult:  MatchResultSuccess,
			wantMatched: true,
		},
		{
			name: "Simplified match: match place",
			matcher: database.MatcherRuntime{
				Matcher:        &goserver.Matcher{Simplified: true},
				Keywords:       []string{"Prague"},
				KeywordOutputs: []string{"Prague"},
				KeywordRegexps: []*regexp.Regexp{
					regexp.MustCompile(`(?i)\bPrague\b`),
				},
			},
			transaction: goserver.Transaction{
				Place: "Prague Airport",
			},
			wantResult:  MatchResultSuccess,
			wantMatched: true,
		},
		{
			name: "Simplified match: match partner name",
			matcher: database.MatcherRuntime{
				Matcher:        &goserver.Matcher{Simplified: true},
				Keywords:       []string{"Lidl"},
				KeywordOutputs: []string{"Lidl"},
				KeywordRegexps: []*regexp.Regexp{
					regexp.MustCompile(`(?i)\bLidl\b`),
				},
			},
			transaction: goserver.Transaction{
				PartnerName: "Lidl Ceska Republika",
			},
			wantResult:  MatchResultSuccess,
			wantMatched: true,
		},
		{
			name: "Simplified match: match|output format",
			matcher: database.MatcherRuntime{
				Matcher:        &goserver.Matcher{Simplified: true},
				Keywords:       []string{"Uber"},
				KeywordOutputs: []string{"Uber Rides"},
				KeywordRegexps: []*regexp.Regexp{
					regexp.MustCompile(`(?i)\bUber\b`),
				},
			},
			transaction: goserver.Transaction{
				Description: "Uber Eats",
			},
			wantResult:  MatchResultSuccess,
			wantMatched: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchWithDetails(&tt.matcher, &tt.transaction)
			assert.Equal(t, tt.wantResult, got.Result)
			assert.Equal(t, tt.wantMatched, got.Matched)
			if tt.name == "Simplified match: match|output format" {
				assert.Equal(t, "Uber", got.MatchedKeyword)
				assert.Equal(t, "Uber Rides", got.MatchedOutput)
			}
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
				PlaceRegexp:       regexp.MustCompile(`(?i)\bAliExpress\b`),
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
