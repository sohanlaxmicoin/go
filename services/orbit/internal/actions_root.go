package orbit

import (
	"github.com/rover/go/services/orbit/internal/ledger"
	"github.com/rover/go/services/orbit/internal/resource"
	"github.com/rover/go/support/render/hal"
)

// RootAction provides a summary of the orbit instance and links to various
// useful endpoints
type RootAction struct {
	Action
}

// JSON renders the json response for RootAction
func (action *RootAction) JSON() {
	var res resource.Root
	res.Populate(
		action.Ctx,
		ledger.CurrentState(),
		action.App.horizonVersion,
		action.App.coreVersion,
		action.App.networkPassphrase,
		action.App.protocolVersion,
	)

	hal.Render(action.W, res)
}
