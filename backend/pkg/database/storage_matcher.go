package database

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

func (s *storage) GetMatchers(userID string) ([]goserver.Matcher, error) {
	result, err := s.db.Model(&models.Matcher{}).Where("user_id = ?", userID).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	matchers := make([]goserver.Matcher, 0)
	for result.Next() {
		var m models.Matcher
		if err := s.db.ScanRows(result, &m); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		matchers = append(matchers, m.FromDB())
	}

	return matchers, nil
}

func (s *storage) CreateMatcher(userID string, matcher goserver.MatcherNoIdInterface) (goserver.Matcher, error) {
	data := models.MatcherToDB(matcher, userID)
	data.ID = uuid.New()
	if err := s.db.Create(data).Error; err != nil {
		return goserver.Matcher{}, fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, userID, "Matcher", data.ID.String(), "CREATED", nil, data); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	s.log.Info("Matcher created", "id", data.ID)

	return data.FromDB(), nil
}

func (s *storage) createMatcherRuntime(m goserver.Matcher) (MatcherRuntime, error) {
	runtime := MatcherRuntime{Matcher: &m}
	if m.DescriptionRegExp != "" {
		r, err := regexp.Compile(m.DescriptionRegExp)
		if err != nil {
			return MatcherRuntime{}, fmt.Errorf("failed to compile description regexp: %w", err)
		}
		runtime.DescriptionRegexp = r
	}

	if m.PartnerAccountNumberRegExp != "" {
		r, err := regexp.Compile(m.PartnerAccountNumberRegExp)
		if err != nil {
			return MatcherRuntime{}, fmt.Errorf("failed to compile partner account regexp: %w", err)
		}
		runtime.PartnerAccountRegexp = r
	}

	if m.PartnerNameRegExp != "" {
		r, err := regexp.Compile(m.PartnerNameRegExp)
		if err != nil {
			return MatcherRuntime{}, fmt.Errorf("failed to compile partner name regexp: %w", err)
		}
		runtime.PartnerNameRegexp = r
	}

	if m.CurrencyRegExp != "" {
		r, err := regexp.Compile(m.CurrencyRegExp)
		if err != nil {
			return MatcherRuntime{}, fmt.Errorf("failed to compile currency regexp: %w", err)
		}
		runtime.CurrencyRegexp = r
	}

	if m.PlaceRegExp != "" {
		r, err := regexp.Compile(m.PlaceRegExp)
		if err != nil {
			return MatcherRuntime{}, fmt.Errorf("failed to compile place regexp: %w", err)
		}
		runtime.PlaceRegexp = r
	}

	if m.Simplified {
		runtime.Keywords = make([]string, len(m.Keywords))
		runtime.KeywordOutputs = make([]string, len(m.Keywords))
		runtime.KeywordRegexps = make([]*regexp.Regexp, len(m.Keywords))
		for i, k := range m.Keywords {
			matcherPart := k
			outputPart := k

			if idx := strings.Index(k, "|"); idx != -1 {
				matcherPart = k[:idx]
				outputPart = k[idx+1:]
			}

			runtime.Keywords[i] = matcherPart
			runtime.KeywordOutputs[i] = outputPart

			// Case-insensitive, whole-word matching
			// We wrap the keyword in \b (word boundary)
			r, err := regexp.Compile(`(?i)\b` + regexp.QuoteMeta(matcherPart) + `\b`)
			if err != nil {
				return MatcherRuntime{}, fmt.Errorf("failed to compile keyword regexp %q: %w", matcherPart, err)
			}
			runtime.KeywordRegexps[i] = r
		}
	}
	return runtime, nil
}

// CreateMatcherRuntimeFromNoId creates a MatcherRuntime from a MatcherNoId (without needing to save to DB first).
// This is useful for testing matchers before they are saved.
//
//nolint:stylecheck
func (s *storage) CreateMatcherRuntimeFromNoId(m goserver.MatcherNoIdInterface) (MatcherRuntime, error) {
	// Convert MatcherNoId to Matcher by creating a temporary matcher with empty ID
	matcher := goserver.Matcher{
		OutputDescription:          m.GetOutputDescription(),
		OutputAccountId:            m.GetOutputAccountId(),
		OutputTags:                 m.GetOutputTags(),
		CurrencyRegExp:             m.GetCurrencyRegExp(),
		PartnerNameRegExp:          m.GetPartnerNameRegExp(),
		PartnerAccountNumberRegExp: m.GetPartnerAccountNumberRegExp(),
		DescriptionRegExp:          m.GetDescriptionRegExp(),
		ExtraRegExp:                m.GetExtraRegExp(),
		PlaceRegExp:                m.GetPlaceRegExp(),
		ConfirmationHistory:        m.GetConfirmationHistory(),
		Simplified:                 m.GetSimplified(),
		Keywords:                   m.GetKeywords(),
	}

	return s.createMatcherRuntime(matcher)
}

func (s *storage) GetMatcherRuntime(userID, id string) (MatcherRuntime, error) {
	m, err := s.GetMatcher(userID, id)
	if err != nil {
		return MatcherRuntime{}, err
	}

	return s.createMatcherRuntime(m)
}

func (s *storage) GetMatchersRuntime(userID string) ([]MatcherRuntime, error) {
	matchers, err := s.GetMatchers(userID)
	if err != nil {
		return nil, err
	}

	res := make([]MatcherRuntime, 0, len(matchers))
	for _, m := range matchers {
		runtime, err := s.createMatcherRuntime(m)
		if err != nil {
			return nil, err
		}

		res = append(res, runtime)
	}

	return res, nil
}

func (s *storage) UpdateMatcher(userID string, id string, matcher goserver.MatcherNoIdInterface,
) (goserver.Matcher, error) {
	return performUpdate[models.Matcher, goserver.MatcherNoIdInterface, goserver.Matcher](s, userID, "Matcher", id, matcher,
		models.MatcherToDB,
		func(m *models.Matcher) goserver.Matcher { return m.FromDB() },
		func(m *models.Matcher, id uuid.UUID) { m.ID = id },
	)
}

func (s *storage) GetMatcher(userID string, id string) (goserver.Matcher, error) {
	var data models.Matcher
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Matcher{}, ErrNotFound
		}

		return goserver.Matcher{}, fmt.Errorf(StorageError, err)
	}

	return data.FromDB(), nil
}

func (s *storage) DeleteMatcher(userID string, id string) error {
	var data models.Matcher
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err == nil {
		if err := s.recordAuditLog(s.db, userID, "Matcher", id, "DELETED", &data, nil); err != nil {
			s.log.Error("Failed to record audit log", "error", err)
		}
	}

	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Matcher{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

// AddMatcherConfirmation atomically appends a confirmation boolean to the matcher's
// confirmation history and trims it to the configured maximum length.
func (s *storage) AddMatcherConfirmation(userID string, id string, confirmed bool) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var m models.Matcher
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&m).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.log.Warn("Matcher not found when adding confirmation", "userID", userID, "matcherID", id)
				return ErrNotFound
			}
			s.log.Error("DB error when loading matcher for confirmation", "error", err)
			return fmt.Errorf(StorageError, err)
		}

		// Use the model helper to add confirmation and respect config max length
		m.AddConfirmation(confirmed, s.cfg.MatcherConfirmationHistoryMax)

		if err := tx.Save(&m).Error; err != nil {
			s.log.Error("DB error when saving matcher after adding confirmation", "error", err)
			return fmt.Errorf(StorageError, err)
		}

		return nil
	})
}
