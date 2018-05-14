package test

import (
	"github.com/stellar/go/services/orbit/internal/test/scenarios"
)

func loadScenario(scenarioName string, includeHorizon bool) {
	stellarCorePath := scenarioName + "-core.sql"
	horizonPath := scenarioName + "-orbit.sql"

	if !includeHorizon {
		horizonPath = "blank-orbit.sql"
	}

	scenarios.Load(StellarCoreDatabaseURL(), stellarCorePath)
	scenarios.Load(DatabaseURL(), horizonPath)
}
