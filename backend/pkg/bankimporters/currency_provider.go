package bankimporters

import (
	"context"
	"fmt"
	"sync"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type DefaultCurrencyProvider struct {
	db     database.Storage
	userID string

	mu         sync.RWMutex
	currencies map[string]string // name -> id
}

func NewDefaultCurrencyProvider(db database.Storage, userID string, initialCurrencies []goserver.Currency) *DefaultCurrencyProvider {
	cache := make(map[string]string)
	for _, c := range initialCurrencies {
		cache[c.Name] = c.Id
	}
	return &DefaultCurrencyProvider{
		db:         db,
		userID:     userID,
		currencies: cache,
	}
}

func (p *DefaultCurrencyProvider) GetCurrencyIdByName(_ context.Context, name string) (string, error) {
	p.mu.RLock()
	id, ok := p.currencies[name]
	p.mu.RUnlock()
	if ok {
		return id, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double check after acquiring lock
	if id, ok := p.currencies[name]; ok {
		return id, nil
	}

	// Not in cache, try to fetch all currencies from DB (maybe created by another process/thread)
	currencies, err := p.db.GetCurrencies(p.userID)
	if err != nil {
		return "", fmt.Errorf("can't get currencies from DB: %w", err)
	}

	// Clear cache, since currencies may have deleted
	p.currencies = make(map[string]string)

	// Fill cache
	for _, c := range currencies {
		p.currencies[c.Name] = c.Id
		if c.Name == name {
			return c.Id, nil
		}
	}

	// Still not found, create it
	newCur, err := p.db.CreateCurrency(p.userID, &goserver.CurrencyNoId{
		Name:        name,
		Description: "Automatically created during bank import",
	})
	if err != nil {
		return "", fmt.Errorf("can't create currency %q: %w", name, err)
	}

	p.currencies[newCur.Name] = newCur.Id
	return newCur.Id, nil
}

type SimpleCurrencyProvider struct {
	currencies map[string]string // name -> id
}

func NewSimpleCurrencyProvider(initialCurrencies []goserver.Currency) *SimpleCurrencyProvider {
	cache := make(map[string]string)
	for _, c := range initialCurrencies {
		cache[c.Name] = c.Id
	}
	return &SimpleCurrencyProvider{
		currencies: cache,
	}
}

func (p *SimpleCurrencyProvider) GetCurrencyIdByName(_ context.Context, name string) (string, error) {
	if id, ok := p.currencies[name]; ok {
		return id, nil
	}
	return "", fmt.Errorf("currency %q not found in simple provider", name)
}
