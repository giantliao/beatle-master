package api

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-master/wallet"
	"github.com/giantliao/beatles-protocol/licenses"
	"github.com/giantliao/beatles-protocol/meta"
	"github.com/giantliao/beatles-protocol/miners"
	"github.com/kprc/libeth/account"
	w2 "github.com/kprc/libeth/wallet"
	"github.com/kprc/nbsnetwork/tools"
	"io/ioutil"
	"net/http"
)

type ListMiners struct {
}

func IsValidLicense(cid account.BeatleAddress, w w2.WalletIntf, l *licenses.License) bool {
	if l.Content.Receiver != cid {
		return false
	}

	if w.BtlAddress() != l.Content.Provider {
		return false
	}

	now := tools.GetNowMsTime()
	if l.Content.ExpireTime < now {
		return false
	}

	forsig, _ := json.Marshal(l.Content)

	bsig := w.BtlSign(forsig)

	ssig := base58.Encode(bsig)

	if ssig != l.Signature {
		return false
	}
	return true
}

func getBestMiners() *miners.BestMiners {
	ms := &miners.BestMiners{}

	mdb := db.GetMinersDb()

	mdb.Iterator()
	for {
		id, md, err := mdb.Next()
		if err != nil || id == "" {
			break
		}

		m := miners.Miner{}
		m.MinerId = md.ID
		m.Location = md.Location
		m.Ipv4Addr = md.Ipv4Addr
		m.Port = md.Port

		ms.Ms = append(ms.Ms, m)

	}
	return ms
}

func (lm *ListMiners) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	var body []byte
	var err error

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	req := &meta.Meta{ContentS: string(body)}

	var (
		sender    string
		cipherTxt []byte
	)
	sender, cipherTxt, err = req.UnMarshal()
	if err != nil || !(account.BeatleAddress(sender).IsValid()) {
		w.WriteHeader(500)
		fmt.Fprintf(w, "not a correct request")
		return
	}
	var wal w2.WalletIntf
	wal, err = wallet.GetWallet()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "server have no wallet")
		return
	}
	var key []byte
	key, err = wal.AesKey2(account.BeatleAddress(sender))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	l := &licenses.License{}
	err = l.UnMarshal(key, cipherTxt)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	if !IsValidLicense(account.BeatleAddress(sender), wal, l) {
		w.WriteHeader(500)
		fmt.Fprintf(w, "signature error")
		return
	}

	ms := getBestMiners()

	cipherTxt, err = ms.Marshal(key)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	resp := &meta.Meta{}
	resp.Marshal(wal.BtlAddress().String(), cipherTxt)

	w.WriteHeader(200)
	fmt.Fprint(w, resp.Content)

	return

}
