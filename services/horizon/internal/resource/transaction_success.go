package resource

import (
	"context"

	"github.com/rover/go/services/horizon/internal/httpx"
	"github.com/rover/go/services/horizon/internal/render/hal"
	"github.com/rover/go/services/horizon/internal/txsub"
)

// Populate fills out the details
func (res *TransactionSuccess) Populate(ctx context.Context, result txsub.Result) {
	res.Hash = result.Hash
	res.Ledger = result.LedgerSequence
	res.Env = result.EnvelopeXDR
	res.Result = result.ResultXDR
	res.Meta = result.ResultMetaXDR

	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	res.Links.Transaction = lb.Link("/transactions", result.Hash)
	return
}
