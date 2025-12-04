package bankimporters

import (
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type Importer interface {
	// Import transactions from the source and returns them.
	// Import(ctx context.Context) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error)

	// Parse and import transactions from the source and returns them. Format could be for example 'csv', 'xslx', etc.
	ParseAndImport(format, data string) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error)
}
