package internal

import (
	"net/http"
	"net/url"

	"github.com/rover/go/clients/orbit"
	"github.com/rover/go/strkey"
	"github.com/rover/go/support/render/hal"
	"github.com/rover/go/support/render/problem"
)

// DakibotHandler causes an account at `Address` to be created.
type DakibotHandler struct {
	Dakibot *Bot
}

// Handle is a method that implements http.HandlerFunc
func (handler *DakibotHandler) Handle(w http.ResponseWriter, r *http.Request) {
	result, err := handler.doHandle(r)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	hal.Render(w, *result)
}

// doHandle is just a convenience method that returns the object to be rendered
func (handler *DakibotHandler) doHandle(r *http.Request) (*orbit.TransactionSuccess, error) {
	err := handler.checkEnabled()
	if err != nil {
		return nil, err
	}

	err = r.ParseForm()
	if err != nil {
		return nil, err
	}

	address, err := handler.loadAddress(r)
	if err != nil {
		return nil, problem.MakeInvalidFieldProblem("addr", err)
	}

	return handler.loadResult(address)
}

func (handler *DakibotHandler) checkEnabled() error {
	if handler.Dakibot != nil {
		return nil
	}

	return &problem.P{
		Type:   "dakibot_disabled",
		Title:  "Dakibot is disabled",
		Status: http.StatusForbidden,
		Detail: "Dakibot is disabled on this network. Contact the server administrator if you believe this to be in error.",
	}
}

func (handler *DakibotHandler) loadAddress(r *http.Request) (string, error) {
	address := r.Form.Get("addr")
	unescaped, err := url.QueryUnescape(address)
	if err != nil {
		return unescaped, err
	}

	_, err = strkey.Decode(strkey.VersionByteAccountID, unescaped)
	return unescaped, err
}

func (handler *DakibotHandler) loadResult(address string) (*orbit.TransactionSuccess, error) {
	result, err := handler.Dakibot.Pay(address)
	switch e := err.(type) {
	case orbit.Error:
		return result, e.Problem.ToProblem()
	case *orbit.Error:
		return result, e.Problem.ToProblem()
	}
	return result, err
}
