package common

import "context"

type ForcedImportKeyType string

const ForcedImportKey ForcedImportKeyType = "forced_import_channel"

type ForcedImport struct {
	UserID         string
	BankImporterID string
}

func GetForcedImportChannel(ctx context.Context) chan<- ForcedImport {
	if ctx == nil {
		return nil
	}
	res, _ := ctx.Value(ForcedImportKey).(chan<- ForcedImport)
	return res
}
