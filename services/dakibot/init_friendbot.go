package main

import (
	"net/http"

	"github.com/rover/go/clients/horizon"
	"github.com/rover/go/services/dakibot/internal"
	"github.com/rover/go/strkey"
)

func initDakibot(dakibotSecret string, networkPassphrase string, horizonURL string, startingBalance string) *internal.Bot {
	if dakibotSecret == "" || networkPassphrase == "" || horizonURL == "" || startingBalance == "" {
		return nil
	}

	// ensure its a seed if its not blank
	strkey.MustDecode(strkey.VersionByteSeed, dakibotSecret)

	return &internal.Bot{
		Secret: dakibotSecret,
		Horizon: &horizon.Client{
			URL:  horizonURL,
			HTTP: http.DefaultClient,
		},
		Network:         networkPassphrase,
		StartingBalance: startingBalance,
	}
}
