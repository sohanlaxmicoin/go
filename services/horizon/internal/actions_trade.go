package horizon

import (
	"errors"

	"fmt"

	"github.com/stellar/go/xdr"
	"github.com/stellar/go/services/horizon/internal/db2"
	"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/services/horizon/internal/render/hal"
	"github.com/stellar/go/services/horizon/internal/resource"
)

type TradeIndexAction struct {
	Action
	OfferFilter        int64
	BaseAssetFilter    xdr.Asset
	CounterAssetFilter xdr.Asset
	PagingParams       db2.PageQuery
	Records            []history.Trade
	Page               hal.Page
}

// JSON is a method for actions.JSON
func (action *TradeIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func flip(a int64, b int64) (int64, int64, bool) {
	if a > b {
		return b, a, true
	}
	return a, b, false
}

// loadParams sets action.Query from the request params
func (action *TradeIndexAction) loadParams() {
	action.PagingParams = action.GetPageQuery()
	action.OfferFilter = action.GetInt64("offer_id")
	action.BaseAssetFilter = action.MaybeGetAsset("base_")
	action.CounterAssetFilter = action.MaybeGetAsset("counter_")
}

// Ensure base is aligned with user's requirement
func ensureTradeDirection(t *history.Trade, requiredBaseId int64) {
	if t.BaseAssetId != requiredBaseId {
		t.BaseAssetId, t.CounterAssetId = t.CounterAssetId, t.BaseAssetId
		t.BaseVolume, t.CounterVolume = t.CounterVolume, t.BaseVolume
	}
}

// loadRecords populates action.Records
func (action *TradeIndexAction) loadRecords() {
	trades := action.HistoryQ().Trades()

	if action.OfferFilter > int64(0) {
		trades = trades.ForOffer(action.OfferFilter)
	}

	baseAssetId, counterAssetId := int64(-1), int64(-1)
	if (action.BaseAssetFilter != xdr.Asset{}) {
		baseAssetId, action.Err = action.HistoryQ().GetAssetID(action.BaseAssetFilter)
		if action.Err != nil {
			return
		}
		if (action.CounterAssetFilter != xdr.Asset{}) {
			counterAssetId, action.Err = action.HistoryQ().GetAssetID(action.CounterAssetFilter)
			if action.Err != nil {
				return
			}
			trades.ForAssetPair(baseAssetId, counterAssetId)
		} else {
			trades.ForSingleAsset(baseAssetId)
		}
	}

	action.Err = trades.Page(action.PagingParams).Select(&action.Records)
	if action.Err != nil {
		return
	}

	// if a base was specified, ensure directionality according to user's requirement
	if baseAssetId != -1 {
		for _, record := range action.Records {
			ensureTradeDirection(&record, baseAssetId)
		}
	}
}

// loadPage populates action.Page
func (action *TradeIndexAction) loadPage() {
	for _, record := range action.Records {
		var res resource.Trade

		action.Err = res.Populate(action.Ctx, record)
		if action.Err != nil {
			return
		}

		action.Page.Add(res)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}

type TradeEffectIndexAction struct {
	Action
	AccountFilter string
	PagingParams  db2.PageQuery
	Records       []history.Effect
	Ledgers       history.LedgerCache
	Page          hal.Page
}

// JSON is a method for actions.JSON
func (action *TradeEffectIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecords,
		action.loadLedgers,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

// loadLedgers populates the ledger cache for this action
func (action *TradeEffectIndexAction) loadLedgers() {
	if action.Err != nil {
		return
	}

	for _, trade := range action.Records {
		action.Ledgers.Queue(trade.LedgerSequence())
	}

	action.Err = action.Ledgers.Load(action.HistoryQ())
}

func (action *TradeEffectIndexAction) loadParams() {
	action.AccountFilter = action.GetString("account_id")
	action.PagingParams = action.GetPageQuery()
}

func (action *TradeEffectIndexAction) loadRecords() {
	trades := action.HistoryQ().Effects().OfType(history.EffectTrade).ForAccount(action.AccountFilter)
	action.Err = trades.Page(action.PagingParams).Select(&action.Records)
}

// loadPage populates action.Page
func (action *TradeEffectIndexAction) loadPage() {
	for _, record := range action.Records {
		var res resource.TradeFromEffect

		ledger, found := action.Ledgers.Records[record.LedgerSequence()]
		if !found {
			msg := fmt.Sprintf("could not find ledger data for sequence %d", record.LedgerSequence())
			action.Err = errors.New(msg)
			return
		}

		action.Err = res.PopulateFromEffect(action.Ctx, record, ledger)
		if action.Err != nil {
			return
		}

		action.Page.Add(res)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
