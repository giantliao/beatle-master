package api

import (
	"fmt"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-protocol/licenses"
	"github.com/giantliao/beatles-protocol/meta"
	"github.com/kprc/libeth/account"
	"net/http"
)

type FreshLicenseSrv struct {

}

func (fls *FreshLicenseSrv)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key, cipherTxt, sender, wal, err := DecodeMeta(r)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	req:=&licenses.FreshLicenseReq{}
	err = req.UnMarshal(key, cipherTxt)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	if sender != req.Receiver.String(){
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	flr:=getLicenseFromDB(req.Receiver)
	if flr == nil{
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
	}

	cipherTxt, err = flr.Marshal(key)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	resp := &meta.Meta{}
	resp.Marshal(wal.BtlAddress().String(), cipherTxt)

	w.WriteHeader(200)
	fmt.Fprint(w, resp.ContentS)

	return

}

func getLicenseFromDB(reciver account.BeatleAddress) *licenses.FreshLicensResult  {

	ld:=db.GetLicenseDb().Find(reciver)
	if ld == nil{
		return nil
	}

	l:=&licenses.FreshLicensResult{}
	l.TxStr = ld.LastTx

	c:=&l.Content
	l.Signature = ld.Sig
	c.Email = ld.Email
	c.Name = ld.Name
	c.Cell = ld.Cell
	c.ExpireTime = ld.ExpireTime
	c.Provider = ld.ServerId
	c.Receiver = ld.CID

	return l


}