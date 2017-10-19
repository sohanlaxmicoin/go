package history

import (
	"fmt"
	"math"

	sq "github.com/Masterminds/squirrel"
	"github.com/stellar/go/services/horizon/internal/db2"
)

// PagingToken returns a cursor for this trade
func (r *Trade) PagingToken() string {
	return fmt.Sprintf("%d-%d", r.HistoryOperationID, r.Order)
}

// Trades provides a helper to filter rows from the `history_trades` table
// with pre-defined filters.  See `TradesQ` methods for the available filters.
func (q *Q) Trades() *TradesQ {
	return &TradesQ{
		parent: q,
		sql:    selectTrade,
	}
}

// ForOffer filters the trade query to only return trades that occurred against
// the offer identified by `id`.
func (q *TradesQ) ForOffer(id int64) *TradesQ {
	q.sql = q.sql.Where("offer_id = ?", id)
	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *TradesQ) Page(page db2.PageQuery) *TradesQ {
	if q.Err != nil {
		return q
	}

	op, idx, err := page.CursorInt64Pair(db2.DefaultPairSep)
	if err != nil {
		q.Err = err
		return q
	}

	// constrain the second portion of the cursor pair to 32-bits
	if idx > math.MaxInt32 {
		idx = math.MaxInt32
	}

	switch page.Order {
	case "asc":
		q.sql = q.sql.
			Where(`(
					 htrd.history_operation_id > ?
				OR (
							htrd.history_operation_id = ?
					AND htrd.order > ?
				))`, op, op, idx).
			OrderBy("htrd.history_operation_id asc, htrd.order asc")
	case "desc":
		q.sql = q.sql.
			Where(`(
					 htrd.history_operation_id < ?
				OR (
							htrd.history_operation_id = ?
					AND htrd.order < ?
				))`, op, op, idx).
			OrderBy("htrd.history_operation_id desc, htrd.order desc")
	}

	q.sql = q.sql.Limit(page.Limit)
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *TradesQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

var selectTrade = sq.Select(
	"history_operation_id",
	"\"order\"",
	"ledger_closed_at",
	"offer_id",
	"base_asset_id",
	"base.asset_type as base_asset_type",
	"base.asset_code as base_asset_code",
	"base.asset_issuer as base_asset_issuer",
	"base_volume",
	"counter_asset_id",
	"counter.asset_type as counter_asset_type",
	"counter.asset_code as counter_asset_code",
	"counter.asset_issuer as counter_asset_issuer",
	"counter_volume",
	"base_is_seller",
).From("history_trades htrd").
	Join("history_assets base ON base_asset_id = base.id").
	Join("history_assets counter ON counter_asset_id = counter.id")

// Filters query for a specific asset id appearing as base or counter
func (q *TradesQ) ForSingleAsset(assetId int64) *TradesQ {
	q.sql = q.sql.Where(sq.Or{sq.Eq{"base_asset_id": assetId}, sq.Eq{"counter_asset_id": assetId}})
	return q
}

// Filters query for a specific asset pair
func (q *TradesQ) ForAssetPair(baseAssetId int64, counterAssetId int64) *TradesQ {
	if baseAssetId>counterAssetId {
		baseAssetId, counterAssetId = counterAssetId, baseAssetId
	}
	q.sql = q.sql.Where(sq.Eq{"base_asset_id": baseAssetId, "counter_asset_id": counterAssetId})
	return q
}