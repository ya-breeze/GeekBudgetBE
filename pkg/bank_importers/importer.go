package bank_importers

import (
	"context"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type Importer interface {
	// Import imports transactions from the source and returns them.
	Import(ctx context.Context) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error)
}
