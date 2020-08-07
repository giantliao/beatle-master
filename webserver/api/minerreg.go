package api

import (
	"fmt"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-protocol/miners"

	"net/http"
)

type MinerRegister struct {
}

func (mr *MinerRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key, cipherTxt, _, _, err := DecodeMeta(r)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	m:=miners.Miner{}
	err = m.UnMarshal(key,cipherTxt)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	mdb:=db.GetMinersDb()
	if err = mdb.Insert(m.Ipv4Addr,m.Port,m.Location,m.MinerId); err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("success"))

	return
}
