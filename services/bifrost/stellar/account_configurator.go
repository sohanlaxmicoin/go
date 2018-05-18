package stellar

import (
	"net/http"
	"time"

	"github.com/rover/go/clients/orbit"
	"github.com/rover/go/keypair"
	"github.com/rover/go/services/bifrost/common"
	"github.com/rover/go/support/errors"
	"github.com/rover/go/support/log"
)

func (ac *AccountConfigurator) Start() error {
	ac.log = common.CreateLogger("StellarAccountConfigurator")
	ac.log.Info("RoverAccountConfigurator starting")

	kp, err := keypair.Parse(ac.IssuerPublicKey)
	if err != nil || (err == nil && ac.IssuerPublicKey[0] != 'G') {
		err = errors.Wrap(err, "Invalid IssuerPublicKey")
		ac.log.Error(err)
		return err
	}

	kp, err = keypair.Parse(ac.SignerSecretKey)
	if err != nil || (err == nil && ac.SignerSecretKey[0] != 'S') {
		err = errors.Wrap(err, "Invalid SignerSecretKey")
		ac.log.Error(err)
		return err
	}

	ac.signerPublicKey = kp.Address()

	root, err := ac.Orbit.Root()
	if err != nil {
		err = errors.Wrap(err, "Error loading Orbit root")
		ac.log.Error(err)
		return err
	}

	if root.NetworkPassphrase != ac.NetworkPassphrase {
		return errors.Errorf("Invalid network passphrase (have=%s, want=%s)", root.NetworkPassphrase, ac.NetworkPassphrase)
	}

	err = ac.updateSequence()
	if err != nil {
		err = errors.Wrap(err, "Error loading issuer sequence number")
		ac.log.Error(err)
		return err
	}

	go ac.logStats()
	return nil
}

func (ac *AccountConfigurator) logStats() {
	for {
		ac.log.WithField("currently_processing", ac.processingCount).Info("Stats")
		time.Sleep(15 * time.Second)
	}
}

// ConfigureAccount configures a new account that participated in ICO.
// * First it creates a new account.
// * Once a trusline exists, it credits it with received number of ETH or BTC.
func (ac *AccountConfigurator) ConfigureAccount(destination, assetCode, amount string) {
	localLog := ac.log.WithFields(log.F{
		"destination": destination,
		"assetCode":   assetCode,
		"amount":      amount,
	})
	localLog.Info("Configuring Rover account")

	ac.processingCountMutex.Lock()
	ac.processingCount++
	ac.processingCountMutex.Unlock()

	defer func() {
		ac.processingCountMutex.Lock()
		ac.processingCount--
		ac.processingCountMutex.Unlock()
	}()

	// Check if account exists. If it is, skip creating it.
	for {
		_, exists, err := ac.getAccount(destination)
		if err != nil {
			localLog.WithField("err", err).Error("Error loading account from Orbit")
			time.Sleep(2 * time.Second)
			continue
		}

		if exists {
			break
		}

		localLog.WithField("destination", destination).Info("Creating Rover account")
		err = ac.createAccount(destination)
		if err != nil {
			localLog.WithField("err", err).Error("Error creating Rover account")
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	if ac.OnAccountCreated != nil {
		ac.OnAccountCreated(destination)
	}

	// Wait for trust line to be created...
	for {
		account, err := ac.Orbit.LoadAccount(destination)
		if err != nil {
			localLog.WithField("err", err).Error("Error loading account to check trustline")
			time.Sleep(2 * time.Second)
			continue
		}

		if ac.trustlineExists(account, assetCode) {
			break
		}

		time.Sleep(2 * time.Second)
	}

	localLog.Info("Trust line found")

	// When trustline found check if needs to authorize, then send token
	if ac.NeedsAuthorize {
		localLog.Info("Authorizing trust line")
		err := ac.allowTrust(destination, assetCode, ac.TokenAssetCode)
		if err != nil {
			localLog.WithField("err", err).Error("Error authorizing trust line")
		}
	}

	localLog.Info("Sending token")
	err := ac.sendToken(destination, assetCode, amount)
	if err != nil {
		localLog.WithField("err", err).Error("Error sending asset to account")
		return
	}

	if ac.OnAccountCredited != nil {
		ac.OnAccountCredited(destination, assetCode, amount)
	}

	localLog.Info("Account successully configured")
}

func (ac *AccountConfigurator) getAccount(account string) (orbit.Account, bool, error) {
	var hAccount orbit.Account
	hAccount, err := ac.Orbit.LoadAccount(account)
	if err != nil {
		if err, ok := err.(*orbit.Error); ok && err.Response.StatusCode == http.StatusNotFound {
			return hAccount, false, nil
		}
		return hAccount, false, err
	}

	return hAccount, true, nil
}

func (ac *AccountConfigurator) trustlineExists(account orbit.Account, assetCode string) bool {
	for _, balance := range account.Balances {
		if balance.Asset.Issuer == ac.IssuerPublicKey && balance.Asset.Code == assetCode {
			return true
		}
	}

	return false
}
