package api

import (
	"errors"
	"github.com/giantliao/beatles-master/wallet"
	"github.com/giantliao/beatles-protocol/meta"
	"github.com/kprc/libeth/account"
	w2 "github.com/kprc/libeth/wallet"
	"io/ioutil"
	"net/http"
)

func DecodeMeta(r *http.Request) (key []byte, cipherTxt []byte, sender string, w w2.WalletIntf, err error) {
	if r.Method != "POST" {
		return nil, nil, "", nil, errors.New("not a post request")
	}
	var body []byte

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return nil, nil, "", nil, errors.New("read http body error")
	}

	req := &meta.Meta{ContentS: string(body)}

	sender, cipherTxt, err = req.UnMarshal()
	if err != nil || !(account.BeatleAddress(sender).IsValid()) {
		return nil, nil, "", nil, errors.New("not a correct request")
	}

	w, err = wallet.GetWallet()
	if err != nil {
		return nil, nil, "", nil, errors.New("server have no wallet")
	}

	key, err = w.AesKey2(account.BeatleAddress(sender))
	if err != nil {
		return nil, nil, "", nil, err
	}

	return
}
