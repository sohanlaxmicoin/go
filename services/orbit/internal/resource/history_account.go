package resource

import (
	"context"

	"github.com/rover/go/services/orbit/internal/db2/history"
)

func (this *HistoryAccount) Populate(ctx context.Context, row history.Account) {
	this.ID = row.Address
	this.AccountID = row.Address
}
