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

	//price.GetEthPrice()
	npsig := &licenses.NoncePriceSig{}
	npc := &npsig.Content

	npc.Nonce = np.Nonce
	npc.Receiver = np.Receiver
	npc.EthAddr = np.EthAddr
	npc.PricePerMonth = config.GetCBtlm().BeatlesPrice
	npc.Month = np.Month
	npc.Total = npc.PricePerMonth * float64(npc.Month)
	npc.EthPrice = price.GetEthPrice()
	if npc.EthPrice == 0 {
		npc.EthPrice = 200
	}
	npc.TotalEth = npc.Total / npc.EthPrice

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
