package wallet

import (
	"errors"
	"github.com/giantliao/beatles-master/config"
	"github.com/kprc/libeth/wallet"
	"github.com/kprc/nbsnetwork/tools"
)

var (
	beatlesMasterWallet wallet.WalletIntf
)

func GetWallet() (wallet.WalletIntf, error) {
	if beatlesMasterWallet == nil {
		return nil, errors.New("no wallet, please load it")
	}
	return beatlesMasterWallet, nil
}

func newWallet(auth, savepath, remoteeth string) wallet.WalletIntf {
	w := wallet.CreateWallet(savepath, remoteeth)
	w.Save(auth)

	return w
}

func LoadWallet(auth string) {
	cfg := config.GetCBtlm()

	if !tools.FileExists(cfg.GetWalletSavePath()) {
		beatlesMasterWallet = newWallet(auth, cfg.GetWalletSavePath(), cfg.EthAccessPoint)
	} else {
		var err error
		beatlesMasterWallet, err = wallet.RecoverWallet(cfg.GetWalletSavePath(), cfg.EthAccessPoint, auth)
		if err != nil {
			panic("load wallet failed : " + err.Error())
		}
	}
}
