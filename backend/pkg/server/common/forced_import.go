package common

import (
	"context"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
)

// Use constants.ForcedImportKey

type ForcedImport struct {
	FamilyID       uuid.UUID
	BankImporterID string
}

func GetForcedImportChannel(ctx context.Context) chan<- ForcedImport {
	if ctx == nil {
		return nil
	}
	res, _ := ctx.Value(constants.ForcedImportKey).(chan<- ForcedImport)
	return res
}
