package common

import (
	"sort"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type txInfo struct {
	tx     goserver.Transaction
	amount decimal.Decimal
}

func AnalyzeDisbalance(targetDelta decimal.Decimal, transactions []goserver.Transaction, accountID, currencyID string) goserver.DisbalanceAnalysis {
	analysis := goserver.DisbalanceAnalysis{
		Delta:            targetDelta,
		TransactionCount: int32(len(transactions)),
		Candidates:       []goserver.DisbalanceCandidate{},
	}

	if len(transactions) == 0 {
		return analysis
	}

	var infos []txInfo
	for _, tx := range transactions {
		amount := decimal.Zero
		for _, m := range tx.Movements {
			if m.AccountId == accountID && m.CurrencyId == currencyID {
				amount = amount.Add(m.Amount)
			}
		}
		if !amount.IsZero() {
			infos = append(infos, txInfo{tx: tx, amount: amount})
		}
	}

	if len(infos) == 0 {
		return analysis
	}

	var candidates []goserver.DisbalanceCandidate
	seenSubsets := make(map[string]bool)

	addCandidate := func(txs []txInfo, candType string) {
		sort.Slice(txs, func(i, j int) bool {
			return txs[i].tx.Id < txs[j].tx.Id
		})

		ids := ""
		candTxs := make([]goserver.DisbalanceCandidateTransaction, len(txs))
		sum := decimal.Zero
		for i, info := range txs {
			ids += info.tx.Id + ","
			candTxs[i] = goserver.DisbalanceCandidateTransaction{
				Id:          info.tx.Id,
				Date:        info.tx.Date,
				Description: info.tx.Description,
				Amount:      info.amount,
			}
			sum = sum.Add(info.amount)
		}

		if seenSubsets[ids] {
			return
		}
		seenSubsets[ids] = true

		diff := sum.Sub(targetDelta).Abs()
		candidates = append(candidates, goserver.DisbalanceCandidate{
			Transactions: candTxs,
			Sum:          sum,
			Difference:   diff,
			Type:         candType,
		})
	}

	// Tier 1: Singles
	for _, info := range infos {
		if info.amount.Equal(targetDelta) {
			addCandidate([]txInfo{info}, "exact_single")
		}
	}

	// Tier 2: Pairs
	for i := 0; i < len(infos); i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[i].amount.Add(infos[j].amount).Equal(targetDelta) {
				addCandidate([]txInfo{infos[i], infos[j]}, "exact_pair")
			}
		}
	}

	// Tier 3: DP Subset Sum (Exact match for 3+ items)
	// We limit N to 50 as per requirements.
	// Since we want multiple subsets and have mixed signs, we use a map-based DP
	// but limit the number of subsets per sum to avoid memory explosion.
	if len(infos) <= 50 {
		targetCents := targetDelta.Mul(decimal.NewFromInt(100)).IntPart()

		// dp[sum] = list of subsets (represented as list of indices)
		dp := make(map[int64][][]int)
		dp[0] = [][]int{{}}

		for idx, info := range infos {
			amountCents := info.amount.Mul(decimal.NewFromInt(100)).IntPart()
			if amountCents == 0 {
				continue
			}

			// Capture current sums to avoid using the same transaction twice in one subset
			currentSums := make([]int64, 0, len(dp))
			for s := range dp {
				currentSums = append(currentSums, s)
			}

			for _, s := range currentSums {
				newSum := s + amountCents
				subsets := dp[s]
				for _, subset := range subsets {
					// Don't add to very large subsets if we already have many candidates
					if len(subset) > 10 && len(candidates) > 20 {
						continue
					}

					newSubset := make([]int, len(subset)+1)
					copy(newSubset, subset)
					newSubset[len(subset)] = idx

					dp[newSum] = append(dp[newSum], newSubset)

					// Limit stored subsets per sum
					if len(dp[newSum]) > 50 {
						break
					}
				}
				// Global safety limit
				if len(dp) > 100000 {
					break
				}
			}
			if len(dp) > 100000 {
				break
			}
		}

		if matchingSubsets, ok := dp[targetCents]; ok {
			for _, subset := range matchingSubsets {
				if len(subset) < 3 {
					// 0, 1, 2 already handled or irrelevant
					continue
				}
				txs := make([]txInfo, len(subset))
				for i, txIdx := range subset {
					txs[i] = infos[txIdx]
				}
				addCandidate(txs, "exact_subset")
			}
		}
	}

	// Rank by subset size (smaller = more likely)
	sort.Slice(candidates, func(i, j int) bool {
		if !candidates[i].Difference.Equal(candidates[j].Difference) {
			return candidates[i].Difference.LessThan(candidates[j].Difference)
		}
		// Smaller subset first
		return len(candidates[i].Transactions) < len(candidates[j].Transactions)
	})

	if len(candidates) > 10 {
		candidates = candidates[:10]
	}

	analysis.Candidates = candidates
	return analysis
}
