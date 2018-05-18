package federation

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/rover/go/clients/orbit"
	"github.com/rover/go/clients/stellartoml"
	"github.com/rover/go/support/http/httptest"
	"github.com/stretchr/testify/assert"
)

func TestLookupByAddress(t *testing.T) {
	hmock := httptest.NewClient()
	tomlmock := &stellartoml.MockClient{}
	c := &Client{StellarTOML: tomlmock, HTTP: hmock}

	// happy path - string integer
	tomlmock.On("GetStellarToml", "rover.network").Return(&stellartoml.Response{
		FederationServer: "https://rover.network/federation",
	}, nil)
	hmock.On("GET", "https://rover.network/federation").
		ReturnJSON(http.StatusOK, map[string]string{
			"rover_address": "scott*rover.network",
			"account_id":      "GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C",
			"memo_type":       "id",
			"memo":            "123",
		})
	resp, err := c.LookupByAddress("scott*rover.network")

	if assert.NoError(t, err) {
		assert.Equal(t, "GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C", resp.AccountID)
		assert.Equal(t, "id", resp.MemoType)
		assert.Equal(t, "123", resp.Memo.String())
	}

	// happy path - integer
	tomlmock.On("GetStellarToml", "rover.network").Return(&stellartoml.Response{
		FederationServer: "https://rover.network/federation",
	}, nil)
	hmock.On("GET", "https://rover.network/federation").
		ReturnJSON(http.StatusOK, map[string]interface{}{
			"rover_address": "scott*rover.network",
			"account_id":      "GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C",
			"memo_type":       "id",
			"memo":            123,
		})
	resp, err = c.LookupByAddress("scott*rover.network")

	if assert.NoError(t, err) {
		assert.Equal(t, "GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C", resp.AccountID)
		assert.Equal(t, "id", resp.MemoType)
		assert.Equal(t, "123", resp.Memo.String())
	}

	// happy path - string
	tomlmock.On("GetStellarToml", "rover.network").Return(&stellartoml.Response{
		FederationServer: "https://rover.network/federation",
	}, nil)
	hmock.On("GET", "https://rover.network/federation").
		ReturnJSON(http.StatusOK, map[string]interface{}{
			"rover_address": "scott*rover.network",
			"account_id":      "GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C",
			"memo_type":       "text",
			"memo":            "testing",
		})
	resp, err = c.LookupByAddress("scott*rover.network")

	if assert.NoError(t, err) {
		assert.Equal(t, "GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C", resp.AccountID)
		assert.Equal(t, "text", resp.MemoType)
		assert.Equal(t, "testing", resp.Memo.String())
	}

	// response exceeds limit
	tomlmock.On("GetStellarToml", "toobig.org").Return(&stellartoml.Response{
		FederationServer: "https://toobig.org/federation",
	}, nil)
	hmock.On("GET", "https://toobig.org/federation").
		ReturnJSON(http.StatusOK, map[string]string{
			"rover_address": strings.Repeat("0", FederationResponseMaxSize) + "*rover.network",
			"account_id":      "GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C",
			"memo_type":       "id",
			"memo":            "123",
		})
	_, err = c.LookupByAddress("response*toobig.org")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "federation response exceeds")
	}

	// failed toml resolution
	tomlmock.On("GetStellarToml", "missing.org").Return(
		(*stellartoml.Response)(nil),
		errors.New("toml failed"),
	)
	resp, err = c.LookupByAddress("scott*missing.org")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "toml failed")
	}

	// 404 federation response
	tomlmock.On("GetStellarToml", "404.org").Return(&stellartoml.Response{
		FederationServer: "https://404.org/federation",
	}, nil)
	hmock.On("GET", "https://404.org/federation").ReturnNotFound()
	resp, err = c.LookupByAddress("scott*404.org")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "failed with (404)")
	}

	// connection error on federation response
	tomlmock.On("GetStellarToml", "error.org").Return(&stellartoml.Response{
		FederationServer: "https://error.org/federation",
	}, nil)
	hmock.On("GET", "https://error.org/federation").ReturnError("kaboom!")
	resp, err = c.LookupByAddress("scott*error.org")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "kaboom!")
	}
}

func TestLookupByID(t *testing.T) {
	horizonMock := &orbit.MockClient{}
	client := &Client{Orbit: horizonMock}

	horizonMock.On("HomeDomainForAccount", "GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C").
		Return("", errors.New("homedomain not set"))

	// an account without a homedomain set fails
	_, err := client.LookupByAccountID("GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C")
	assert.Error(t, err)
	assert.Equal(t, "get homedomain failed: homedomain not set", err.Error())
}

func TestForwardRequest(t *testing.T) {
	hmock := httptest.NewClient()
	tomlmock := &stellartoml.MockClient{}
	c := &Client{StellarTOML: tomlmock, HTTP: hmock}

	// happy path - string integer
	tomlmock.On("GetStellarToml", "rover.network").Return(&stellartoml.Response{
		FederationServer: "https://rover.network/federation",
	}, nil)
	hmock.On("GET", "https://rover.network/federation").
		ReturnJSON(http.StatusOK, map[string]string{
			"account_id": "GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C",
			"memo_type":  "id",
			"memo":       "123",
		})
	fields := url.Values{}
	fields.Add("federation_type", "bank_account")
	fields.Add("swift", "BOPBPHMM")
	fields.Add("acct", "2382376")
	resp, err := c.ForwardRequest("rover.network", fields)

	if assert.NoError(t, err) {
		assert.Equal(t, "GASTNVNLHVR3NFO3QACMHCJT3JUSIV4NBXDHDO4VTPDTNN65W3B2766C", resp.AccountID)
		assert.Equal(t, "id", resp.MemoType)
		assert.Equal(t, "123", resp.Memo.String())
	}
}

func Test_url(t *testing.T) {
	c := &Client{}

	// forward requests
	qstr := url.Values{}
	qstr.Add("type", "forward")
	qstr.Add("federation_type", "bank_account")
	qstr.Add("swift", "BOPBPHMM")
	qstr.Add("acct", "2382376")
	furl := c.url("https://rover.network/federation", qstr)
	assert.Equal(t, "https://rover.network/federation?acct=2382376&federation_type=bank_account&swift=BOPBPHMM&type=forward", furl)

	// regression: ensure that query is properly URI encoded
	qstr = url.Values{}
	qstr.Add("type", "q")
	qstr.Add("q", "scott+receiver1@rover.network*rover.network")
	furl = c.url("", qstr)
	assert.Equal(t, "?q=scott%2Breceiver1%40rover.network%2Arover.network&type=q", furl)
}
