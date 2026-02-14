package common

import (
	"context"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
)

// Use constants.ForcedImportKey

type ForcedImport struct {
	UserID         string
	BankImporterID string
}

func GetForcedImportChannel(ctx context.Context) chan<- ForcedImport {
	if ctx == nil {
		return nil
	}
	res, _ := ctx.Value(constants.ForcedImportKey).(chan<- ForcedImport)
	return res
}
