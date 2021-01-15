package api

import (
	"fmt"
	"github.com/giantliao/beatles-master/config"
	"github.com/giantliao/beatles-master/price"
	"github.com/giantliao/beatles-protocol/licenses"
	"github.com/giantliao/beatles-protocol/meta"
	"net/http"
)

type NoncePriceSrv struct {
}

func (nps *NoncePriceSrv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key, cipherTxt, _, wat, err := DecodeMeta(r)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	np := &licenses.NoncePrice{}
	err = np.UnMarshal(key, cipherTxt)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	var npsig *licenses.NoncePriceSig

	//price.GetEthPrice()
	if np.PayTyp == licenses.PayTypETH {
		npsig = ethNoncePrice(np)

	} else if np.PayTyp == licenses.PayTypBTLC {
		npsig = btlcNoncePrice(np)
	}
	npsig.Content.PayTyp = np.PayTyp

	npsig.Sign(func(data []byte) []byte {
		return wat.BtlSign(data)
	})

	var c []byte

	c, err = npsig.Marshal(key)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	resp := &meta.Meta{}
	resp.Marshal(wat.BtlAddress().String(), c)

	w.WriteHeader(200)
	fmt.Fprintf(w, resp.ContentS)
}

func ethNoncePrice(np *licenses.NoncePrice) *licenses.NoncePriceSig {
	npsig := &licenses.NoncePriceSig{}
	npc := &npsig.Content

	npc.Nonce = np.Nonce
	npc.Receiver = np.Receiver
	npc.Payer = np.Payer
	npc.PricePerMonth = config.GetCBtlm().BeatlesPrice
	npc.Month = np.Month
	npc.Total = npc.PricePerMonth * float64(npc.Month)
	npc.MarketPrice = price.GetEthPrice()
	if npc.MarketPrice == 0 {
		npc.MarketPrice = 200
	}
	npc.TotalPrice = npc.Total / npc.MarketPrice

	return npsig
}

func btlcNoncePrice(np *licenses.NoncePrice) *licenses.NoncePriceSig {
	npsig := &licenses.NoncePriceSig{}
	npc := &npsig.Content

	npc.Nonce = np.Nonce
	npc.Receiver = np.Receiver
	npc.Payer = np.Payer
	npc.PricePerMonth = config.GetCBtlm().BTLCoinPrice
	npc.Month = np.Month
	npc.Total = npc.PricePerMonth * float64(npc.Month)
	npc.MarketPrice = 1.0
	npc.TotalPrice = npc.Total / npc.MarketPrice

	return npsig
}
