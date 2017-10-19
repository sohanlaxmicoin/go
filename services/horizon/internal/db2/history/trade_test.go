package history

import (
	"testing"

	"github.com/stellar/go/build"
	"github.com/stellar/go/services/horizon/internal/db2"
	"github.com/stellar/go/services/horizon/internal/test"
	"github.com/stellar/go/xdr"
)

func TestTradeQueries(t *testing.T) {
	tt := test.Start(t).Scenario("kahuna")
	defer tt.Finish()
	q := &Q{tt.HorizonSession()}

	var trades []Trade

	// All trades
	err := q.Trades().Select(&trades)
	if tt.Assert.NoError(err) {
		tt.Assert.Len(trades, 4)
	}

	// Paging
	pq := db2.MustPageQuery(trades[0].PagingToken(), "asc", 1)
	var pt []Trade

	err = q.Trades().Page(pq).Select(&pt)
	if tt.Assert.NoError(err) {
		tt.Assert.Len(pt, 1)
		tt.Assert.Equal(trades[1], pt[0])
	}

	// Cursor bounds checking
	pq = db2.MustPageQuery("", "desc", 1)
	err = q.Trades().Page(pq).Select(&pt)
	tt.Assert.NoError(err)

	// Test ForOffer
	err = q.Trades().ForOffer(2).Select(&trades)
	if tt.Assert.NoError(err) {
		tt.Assert.Len(trades, 2)
	}

	err = q.Trades().ForOffer(4).Select(&trades)
	if tt.Assert.NoError(err) {
		tt.Assert.Len(trades, 0)
	}

	// test for single asset
	nativeAssetId, _ := q.GetAssetID(build.NativeAsset().MustXDR())
	q.Trades().ForSingleAsset(nativeAssetId).Select(&trades)
	tt.Assert.Len(trades, 4)

	q.Trades().ForSingleAsset(3).Select(&trades)
	tt.Assert.Len(trades, 2)

	// test for asset pairs
	q.Trades().ForAssetPair(2, 3).Select(&trades)
	tt.Assert.Len(trades, 0)

	q.Trades().ForAssetPair(1, 2).Select(&trades)
	tt.Assert.Len(trades, 1)

	tt.Assert.Equal(xdr.Int64(2000000000), trades[0].BaseVolume)
	tt.Assert.Equal(xdr.Int64(1000000000), trades[0].CounterVolume)
	tt.Assert.Equal(true, trades[0].BaseIsSeller)
}
