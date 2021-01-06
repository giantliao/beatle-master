package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-protocol/licenses"
	"github.com/giantliao/beatles-protocol/meta"
	"github.com/kprc/libeth/account"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type PurchaseLicense struct {
	renewLicenseLock sync.Mutex
}

func isEqual(f1, f2 float64) bool {
	d := f1 - f2
	if d < 0 {
		d = f2 - f1
	}

	if d < 0.000001 {
		return true
	}

	return false
}

func (pl *PurchaseLicense) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key, cipherTxt, sender, wal, err := DecodeMeta(r)
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

	if !((&lr.TXSig).ValidSig(wal.BtlAddress().DerivePubKey())) {
		w.WriteHeader(500)
		fmt.Fprintf(w, "signature not correct")
		return
	}

	if lic,err:=findLicense(lr.EthTransaction,lr.TXSig.Content.Receiver);err==nil{
		if err = reply(lic,key,w);err!=nil{
			log.Println("reply license failed",lr.EthTransaction.String())
		}
		return
	}

	log.Println(lr.String())
	log.Println("begin check receipt from block chain...")

	var (
		total float64
		cnt   int
	)

	for {
		total, err = wal.CheckReceiptWithNonce(lr.TXSig.Content.EthAddr, lr.EthTransaction, lr.TXSig.Content.Nonce)
		if err != nil && strings.Contains(err.Error(), "pending") {
			cnt++
			if cnt > 2000 {
				w.WriteHeader(500)
				log.Println(lr.EthTransaction, err.Error())
				fmt.Fprintf(w, err.Error())
				return
			}
			time.Sleep(time.Second)
			fmt.Println("wait for confirm :", lr.EthTransaction.String())
			continue
		} else if err != nil {
			w.WriteHeader(500)
			log.Println(lr.EthTransaction, err.Error())
			fmt.Fprintf(w, err.Error())
			return
		} else {
			break
		}
	}

	if !isEqual(total, lr.TXSig.Content.TotalEth) {
		w.WriteHeader(500)
		errmsg := "eth value not correct"
		log.Println(lr.EthTransaction, errmsg, total, lr.TXSig.Content.TotalEth)
		fmt.Fprintf(w, errmsg)
		return
	}

	rdb := db.GetReceiptsDb()
	err = rdb.Insert(lr.EthTransaction.String(), account.BeatleAddress(sender), "eth",
		lr.TXSig.Content.EthAddr.String(), lr.TXSig.Content.EthPrice, lr.TXSig.Content.Month)
	if err != nil {
		log.Println("receipt :", lr.EthTransaction.String(), " is existed")
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	log.Println("Build a new license for receipt: " + lr.EthTransaction.String())
	//return a new license to client
	pl.renewLicenseLock.Lock()

	ld := db.GetLicenseDb().Find(lr.TXSig.Content.Receiver)
	expireTime := int64(0)
	if ld == nil {
		expireTime = tools.Moth2Expire(0, lr.TXSig.Content.Month)
	} else {
		expireTime = tools.Moth2Expire(ld.ExpireTime, lr.TXSig.Content.Month)
	}
	lc := &licenses.LicenseContent{}
	lc.Provider = wal.BtlAddress()
	lc.Receiver = lr.TXSig.Content.Receiver
	lc.Name = lr.Name
	lc.Email = lr.Email
	lc.Cell = lr.Cell
	lc.ExpireTime = expireTime

	forsig, _ := json.Marshal(*lc)

	bsig := wal.BtlSign(forsig)

	l := &licenses.License{}
	l.Signature = base58.Encode(bsig)
	l.Content = *lc

	err = db.GetLicenseDb().Update(lc.Receiver, lc.Provider, l.Signature, lr.EthTransaction.String(), lc.Name, lc.Email, lc.Cell, lc.ExpireTime)
	if err != nil {
		log.Println("receipt :", lr.EthTransaction.String(), " update to db failed")
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		pl.renewLicenseLock.Unlock()
		return
	}
	pl.renewLicenseLock.Unlock()

	if err = reply(l,key,w);err!=nil{
		log.Println("receipt :", lr.EthTransaction.String(), " marshal license failed")
	}

}

func reply(l *licenses.License, key []byte, w http.ResponseWriter) error  {

	content, err := l.Marshal(key)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
	}

	mresp := &meta.Meta{}
	mresp.Marshal(l.Content.Provider.String(), content)

	w.WriteHeader(200)
	fmt.Fprint(w, mresp.ContentS)

	return nil
}

func findLicense(tx common.Hash, receiver account.BeatleAddress) (*licenses.License,error)  {
	ld := db.GetLicenseDb().Find(receiver)

	if tx.String() != ld.LastTx{
		return nil,errors.New("not found")
	}

	l:=&licenses.License{}

	c:=&l.Content

	l.Signature = ld.Sig
	c.Email = ld.Email
	c.Name = ld.Name
	c.Cell = ld.Cell
	c.ExpireTime = ld.ExpireTime
	c.Provider = ld.ServerId
	c.Receiver = ld.CID

	return l,nil
}