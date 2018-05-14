// Package reap contains the history reaping subsystem for orbit.  This system
// is designed to remove data from the history database such that it does not
// grow indefinitely.  The system can be configured with a number of ledgers to
// maintain at a minimum.
package reap

import (
	"time"

	"github.com/stellar/go/support/db"
)

// System represents the history reaping subsystem of orbit.
type System struct {
	HorizonDB      *db.Session
	RetentionCount uint

	nextRun time.Time
}

// New initializes the reaper, causing it to begin polling the rover-core
// database for now ledgers and ingesting data into the orbit database.
func New(retention uint, orbit *db.Session) *System {
	r := &System{
		HorizonDB:      orbit,
		RetentionCount: retention,
	}

	r.nextRun = time.Now().Add(1 * time.Hour)
	return r
}
