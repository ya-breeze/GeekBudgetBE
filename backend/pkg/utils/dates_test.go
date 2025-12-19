package utils

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = Describe("Dates Utils", func() {
	Context("RoundToGranularity", func() {
		Context("GranularityMonth", func() {
			It("should round down to the first day of the month", func() {
				// 2024-05-15 10:30:00 -> 2024-05-01 00:00:00
				input := time.Date(2024, 5, 15, 10, 30, 0, 0, time.UTC)
				expected := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
				result := RoundToGranularity(input, GranularityMonth, false)
				Expect(result).To(Equal(expected))
			})

			It("should stay at the first day of the month if already there (round down)", func() {
				input := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
				result := RoundToGranularity(input, GranularityMonth, false)
				Expect(result).To(Equal(input))
			})

			It("should round up to the first day of the next month", func() {
				// 2024-05-15 -> 2024-06-01
				input := time.Date(2024, 5, 15, 10, 30, 0, 0, time.UTC)
				expected := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
				result := RoundToGranularity(input, GranularityMonth, true)
				Expect(result).To(Equal(expected))
			})

			It("should stay at the first day of the month if already there (round up - boundary)", func() {
				// Boundary condition: if we are exactly at the start, ceiling should be same
				input := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
				result := RoundToGranularity(input, GranularityMonth, true)
				Expect(result).To(Equal(input))
			})

			It("should handle December correctly when rounding up", func() {
				// 2024-12-15 -> 2025-01-01
				input := time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC)
				expected := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
				result := RoundToGranularity(input, GranularityMonth, true)
				Expect(result).To(Equal(expected))
			})
		})

		Context("GranularityYear", func() {
			It("should round down to the first day of the year", func() {
				// 2024-05-15 -> 2024-01-01
				input := time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)
				expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				result := RoundToGranularity(input, GranularityYear, false)
				Expect(result).To(Equal(expected))
			})

			It("should stay at the first day of the year if already there (round down)", func() {
				input := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				result := RoundToGranularity(input, GranularityYear, false)
				Expect(result).To(Equal(input))
			})

			It("should round up to the first day of the next year", func() {
				// 2024-05-15 -> 2025-01-01
				input := time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)
				expected := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
				result := RoundToGranularity(input, GranularityYear, true)
				Expect(result).To(Equal(expected))
			})

			It("should stay at the first day of the year if already there (round up - boundary)", func() {
				input := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				result := RoundToGranularity(input, GranularityYear, true)
				Expect(result).To(Equal(input))
			})
		})
	})
})
