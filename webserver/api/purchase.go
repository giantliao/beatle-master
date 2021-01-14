package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/giantliao/beatles-master/coin"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-protocol/licenses"
	"github.com/giantliao/beatles-protocol/meta"
	"github.com/kprc/libeth/account"
	"github.com/kprc/libeth/wallet"
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

	if lic,err:=findLicense(lr.TxStr,lr.TXSig.Content.Receiver);err==nil{
		if err = reply(lic,key,w);err!=nil{
			log.Println("reply license failed",lr.TxStr)
		}
		return
	}

	log.Println(lr.String())
	log.Println("begin check receipt from block chain...")

	var total float64

	if lr.TXSig.Content.PayTyp == licenses.PayTypETH{
		total,err = checkEthReceipt(wal,lr)
		if err!=nil{
			w.WriteHeader(500)
			log.Println(lr.TxStr, err.Error())
			fmt.Fprintf(w, err.Error())
			return
		}
	}else{
		var fromaddr common.Address
		var toaddr common.Address
		total,fromaddr,toaddr,err = coin.GetBTLCoinToken().CheckHashAndGet(common.HexToHash(lr.TxStr),2000)
		if err != nil{
			w.WriteHeader(500)
			log.Println(lr.TxStr, err.Error())
			fmt.Fprintf(w, err.Error())
			return
		}

		if fromaddr.String() != lr.TXSig.Content.Payer.String() || toaddr.String() != wal.AccountString(){
			w.WriteHeader(500)
			log.Println(lr.TxStr, "account not correct", fromaddr.String(),toaddr.String())
			fmt.Fprintf(w,"account not correct")

			return
		}
	}

	if !isEqual(total, lr.TXSig.Content.TotalPrice) {
		w.WriteHeader(500)
		errmsg := "eth value not correct"
		log.Println(lr.TxStr, errmsg, total, lr.TXSig.Content.TotalPrice)
		fmt.Fprintf(w, errmsg)
		return
	}

	rdb := db.GetReceiptsDb()
	err = rdb.Insert(lr.TxStr, account.BeatleAddress(sender), "eth",
		lr.TXSig.Content.Payer.String(), lr.TXSig.Content.MarketPrice, lr.TXSig.Content.Month)
	if err != nil {
		log.Println("receipt :", lr.TxStr, " is existed")
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	log.Println("Build a new license for receipt: " + lr.TxStr)
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

	err = db.GetLicenseDb().Update(lc.Receiver, lc.Provider, l.Signature, lr.TxStr, lc.Name, lc.Email, lc.Cell, lc.ExpireTime)
	if err != nil {
		log.Println("receipt :", lr.TxStr, " update to db failed")
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		pl.renewLicenseLock.Unlock()
		return
	}
	pl.renewLicenseLock.Unlock()

	if err = reply(l,key,w);err!=nil{
		log.Println("receipt :", lr.TxStr, " marshal license failed")
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

func findLicense(tx string, receiver account.BeatleAddress) (*licenses.License,error)  {
	ld := db.GetLicenseDb().Find(receiver)

	if ld == nil || tx != ld.LastTx{
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



func checkEthReceipt(wal wallet.WalletIntf, lr *licenses.LicenseRenew) (float64, error)  {
	var (
		total float64
		err error
		cnt   int
	)

	for {
		total, err = wal.CheckReceiptWithNonce(lr.TXSig.Content.Payer, common.HexToHash(lr.TxStr), lr.TXSig.Content.Nonce)
		if err != nil && strings.Contains(err.Error(), "pending") {
			cnt++
			if cnt > 2000 {
				return 0,err
			}
			time.Sleep(time.Second)
			fmt.Println("wait for confirm :", lr.TxStr)
			continue
		} else if err != nil {
			return 0,err
		} else {
			break
		}
	}

	return total,nil
}