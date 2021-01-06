package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/giantliao/beatles-master/config"
	"github.com/kprc/libeth/account"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
	"time"
)

type MinersDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	quit chan struct{}
}

var (
	minersStore     *MinersDb
	minersStoreLock sync.Mutex
)

func newMinersStore() *MinersDb {
	cfg := config.GetCBtlm()

	db := db.NewFileDb(cfg.GetMinersDbFile()).Load()

	return &MinersDb{NbsDbInter: db,quit: make(chan struct{},1)}
}

func GetMinersDb() *MinersDb {
	if minersStore == nil {
		minersStoreLock.Lock()
		defer minersStoreLock.Unlock()

		if minersStore == nil {
			minersStore = newMinersStore()
		}
	}

	return minersStore
}

type MinerDesc struct {
	Ipv4Addr   string                `json:"ipv_4_addr"`
	Port       int                   `json:"port"`
	Location   string                `json:"location"`
	ID         account.BeatleAddress `json:"-"`
	CreateTime int64                 `json:"create_time"`
	UpdateTime int64                 `json:"update_time"`
}

func (md *MinerDesc) String() string {
	msg := ""
	msg += fmt.Sprintf("%-50s", md.ID.String())
	msg += fmt.Sprintf("%-18s", md.Ipv4Addr)
	msg += fmt.Sprintf("%-7d", md.Port)
	msg += fmt.Sprintf("%-12s", md.Location)
	msg += fmt.Sprintf("%-20s", tools.Int64Time2String(md.CreateTime))
	msg += fmt.Sprintf("%-20s", tools.Int64Time2String(md.UpdateTime))

	return msg

}

func (mdb *MinersDb) Insert(ipv4 string, port int, location string, id account.BeatleAddress) error {
	mdb.dbLock.Lock()
	defer mdb.dbLock.Unlock()

	md := MinerDesc{Ipv4Addr: ipv4, Port: port, Location: location, ID: id}

	if _, err := minersStore.NbsDbInter.Find(id.String()); err != nil {
		now := tools.GetNowMsTime()
		md.CreateTime = now
		md.UpdateTime = now
		j, _ := json.Marshal(md)
		minersStore.NbsDbInter.Insert(id.String(), string(j))
		return nil
	} else {
		return errors.New("key is existed, row id is " + id.String())
	}
}

func (mdb *MinersDb) Update(ipv4 string, port int, location string, id account.BeatleAddress) error {
	mdb.dbLock.Lock()
	defer mdb.dbLock.Unlock()

	md := &MinerDesc{}

	now := tools.GetNowMsTime()

	if minerDescStr, err := minersStore.NbsDbInter.Find(id.String()); err != nil {
		md.Ipv4Addr = ipv4
		md.Port = port
		md.Location = location
		md.CreateTime = now
	} else {
		json.Unmarshal([]byte(minerDescStr), md)
		//if md.Ipv4Addr == ipv4 && port == md.Port && location == md.Location {
		//	return errors.New("nothing to update")
		//}

		md.Ipv4Addr = ipv4
		md.Location = location
		md.Port = port
	}

	md.UpdateTime = now

	j, _ := json.Marshal(*md)
	minersStore.NbsDbInter.Update(id.String(), string(j))

	return nil
}

func (mdb *MinersDb) Delete(id account.BeatleAddress) {
	mdb.dbLock.Lock()
	defer mdb.dbLock.Unlock()

	minersStore.NbsDbInter.Delete(id.String())
}

func (mdb *MinersDb) Find(id account.BeatleAddress) (md *MinerDesc, err error) {
	mdb.dbLock.Lock()
	defer mdb.dbLock.Unlock()

	if minerDescStr, err := minersStore.NbsDbInter.Find(id.String()); err != nil {
		return nil, err
	} else {
		md = &MinerDesc{}
		json.Unmarshal([]byte(minerDescStr), md)
		md.ID = id
		return md, nil
	}
}

func (mdb *MinersDb) Save() {
	mdb.dbLock.Lock()
	defer mdb.dbLock.Unlock()

	mdb.NbsDbInter.Save()
}

func (mdb *MinersDb) Iterator() *db.DBCusor {
	mdb.dbLock.Lock()
	defer mdb.dbLock.Unlock()

	return mdb.NbsDbInter.DBIterator()
}

func (mdb *MinersDb) Next(cusor *db.DBCusor) (id account.BeatleAddress, md *MinerDesc, err error) {
	if cusor == nil {
		return account.BeatleAddress(""), nil, errors.New("initialize cursor first")
	}
	mdb.dbLock.Lock()
	k, v := cusor.Next()
	if k == "" {
		mdb.dbLock.Unlock()
		return "", nil, errors.New("no miner in list")
	}
	mdb.dbLock.Unlock()

	md = &MinerDesc{}
	if err := json.Unmarshal([]byte(v), md); err != nil {
		return "", nil, err
	}
	id = account.BeatleAddress(k)
	md.ID = id

	return
}

func (mdb *MinersDb)TimeOut() {
	tic := time.NewTicker(time.Second * 300)
	defer tic.Stop()

	for{
		select {
		case <-tic.C:
			mdb.doTimeOUt()
		case <-mdb.quit:
			return
		}

	}
}

func (mdb *MinersDb)doTimeOUt()  {
	iter:=mdb.Iterator()

	var dKeys []account.BeatleAddress

	now:=tools.GetNowMsTime()

	for{
		id,md,err:=mdb.Next(iter)
		if err!=nil{
			break
		}
		if now - md.UpdateTime > 1800000{
			dKeys = append(dKeys,id)
		}
	}

	for i:=0;i<len(dKeys);i++{
		mdb.Delete(dKeys[i])
	}

}

func (mdb *MinersDb)Close()  {
	close(mdb.quit)
}

func (mdb *MinersDb)StringAll() string  {
	iter := mdb.Iterator()

	msg := ""

	for {
		_, v, err := mdb.Next(iter)
		if err != nil {
			if msg == "" {
				msg = "no miners in db"
			}
			return msg
		}
		if msg != "" {
			msg += "\r\n"
		}
		msg += v.String()

	}

	return msg
}