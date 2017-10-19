package resource

import (
	"errors"
	"fmt"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/services/horizon/internal/httpx"
	"github.com/stellar/go/services/horizon/internal/render/hal"
	"golang.org/x/net/context"
)

// PopulateFromEffect fills out the details of a trade resource from a
// history.Effect row.
func (res *TradeFromEffect) PopulateFromEffect(
	ctx context.Context,
	row history.Effect,
	ledger history.Ledger,
) (err error) {
	if row.Type != history.EffectTrade {
		err = errors.New("invalid effect; not a trade")
		return
	}

	if row.LedgerSequence() != ledger.Sequence {
		err = errors.New("invalid ledger; different sequence than trade")
		return
	}

	row.UnmarshalDetails(res)
	res.ID = row.PagingToken()
	res.PT = row.PagingToken()
	res.Buyer = row.Account
	res.LedgerCloseTime = ledger.ClosedAt
	res.populateLinks(ctx, res.Seller, res.Buyer, row.HistoryOperationID)

	return
}

// Populate fills out the details of a trade using a row from the history_trades
// table.
func (res *Trade) Populate(
	ctx context.Context,
	row history.Trade,
) (err error) {
	res.ID = row.PagingToken()
	res.PT = row.PagingToken()
	res.OfferID = fmt.Sprintf("%d", row.OfferID)
	res.BaseAssetType = row.BaseAssetType
	res.BaseAssetCode = row.BaseAssetCode
	res.BaseAssetIssuer = row.BaseAssetIssuer
	res.BaseVolume = amount.String(row.BaseVolume)
	res.CounterAssetType = row.CounterAssetType
	res.CounterAssetCode = row.CounterAssetCode
	res.CounterAssetIssuer = row.CounterAssetIssuer
	res.CounterVolume = amount.String(row.CounterVolume)
	res.LedgerCloseTime = row.LedgerCloseTime
	res.BaseIsSeller = row.BaseIsSeller
	res.populateLinks(ctx, row.HistoryOperationID)
	return
}

// PagingToken implementation for hal.Pageable
func (res Trade) PagingToken() string {
	return res.PT
}


// PagingToken implementation for hal.Pageable
func (res TradeFromEffect) PagingToken() string {
	return res.PT
}

func (res *TradeFromEffect) populateLinks(
	ctx context.Context,
	seller string,
	buyer string,
	opid int64,
) {
	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	res.Links.Seller = lb.Link("/accounts", res.Seller)
	res.Links.Buyer = lb.Link("/accounts", res.Buyer)
	res.Links.Operation = lb.Link(
		"/operations",
		fmt.Sprintf("%d", opid),
	)
}

func (res *Trade) populateLinks(
	ctx context.Context,
	opid int64,
) {
	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	res.Links.Base = lb.Link("/accounts", res.BaseAssetIssuer)
	res.Links.Counter = lb.Link("/accounts", res.CounterAssetIssuer)
	res.Links.Operation = lb.Link(
		"/operations",
		fmt.Sprintf("%d", opid),
	)
}
