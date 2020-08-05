package api

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/giantliao/beatles-master/config"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-master/wallet"
	"github.com/giantliao/beatles-protocol/licenses"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/giantliao/beatles-protocol/meta"
	"github.com/kprc/libeth/account"
	w2 "github.com/kprc/libeth/wallet"
	"io/ioutil"
	"net/http"
)

type PurchaseLicense struct {
	renewLicenseLock sync.Mutex
}

func (pl *PurchaseLicense) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	lr := &licenses.LicenseRenew{}
	err = lr.UnMarshal(key, cipherTxt)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	log.Println(*lr)
	log.Println("begin check receipt from block chain...")

	var total float64
	cnt := 0
	for {
		total, err = wal.CheckReceipt(lr.EthAddress, lr.EthTransaction)
		if err != nil && strings.Contains(err.Error(), "pending") {
			cnt++
			if cnt > 5 {
				w.WriteHeader(500)
				fmt.Fprintf(w, err.Error())
				return
			}
			continue
		} else if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, err.Error())
			return
		} else {
			break
		}
	}

	cfg := config.GetCBtlm()
	price := lr.CurrentPrice

	if price != cfg.CurrentPrice {
		price = cfg.LastPrice
	}
	if price == 0 {
		log.Println(" price is 0")
		w.WriteHeader(500)
		fmt.Fprintf(w, " price is 0")
		return
	}
	month := int64(total / price)

	rdb := db.GetReceiptsDb()
	err = rdb.Insert(lr.EthTransaction.String(), account.BeatleAddress(sender), "eth", lr.EthAddress.String(), price, month)
	if err != nil {
		log.Println("receipt :", lr.EthTransaction.String(), " is existed")
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	log.Println("receipt: "+lr.EthTransaction.String(),
		"price: ", strconv.FormatFloat(lr.CurrentPrice, 'f', -1, 64),
		"month:", strconv.Itoa(int(lr.Month)),
		"system price:", strconv.FormatFloat(lr.CurrentPrice, 'f', -1, 64),
		"system month:", strconv.Itoa(int(month)))

	//return a new license to client
	pl.renewLicenseLock.Lock()
	defer pl.renewLicenseLock.Unlock()

	ld := db.GetLicenseDb().Find(lr.Receiver)
	expireTime := int64(0)
	if ld == nil {
		expireTime = tools.Moth2Expire(0, month)
	} else {
		expireTime = tools.Moth2Expire(ld.ExpireTime, month)
	}
	lc := &licenses.LicenseContent{}
	lc.Provider = wal.BtlAddress()
	lc.Receiver = lr.Receiver
	lc.Name = lr.Name
	lc.Email = lr.Email
	lc.Cell = lr.Cell
	lc.ExpireTime = expireTime

	forsig, _ := json.Marshal(*lc)

	bsig := wal.BtlSign(forsig)

	l := &licenses.License{}
	l.Signature = base58.Encode(bsig)
	l.Content = *lc

	err = db.GetLicenseDb().Update(lc.Receiver, lc.Provider, l.Signature, lc.Name, lc.Email, lc.Cell, lc.ExpireTime)
	if err != nil {
		log.Println("receipt :", lr.EthTransaction.String(), " update to db failed")
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	var content []byte
	content, err = l.Marshal(key)
	if err != nil {
		log.Println("receipt :", lr.EthTransaction.String(), " marshal license failed")
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
	}

	mresp := &meta.Meta{}
	mresp.Marshal(lc.Provider.String(), content)

	w.WriteHeader(200)
	fmt.Fprint(w, mresp.Content)

	return

}
